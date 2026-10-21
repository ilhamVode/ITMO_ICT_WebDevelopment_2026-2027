package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidUser = errors.New("имя пользователя: 3–50 символов, пароль: 8–72 байта")
var ErrUnauthorized = errors.New("неверное имя пользователя или пароль")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, password string) (User, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 50 || len(password) < 8 || len(password) > 72 {
		return User{}, ErrInvalidUser
	}
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			return User{}, ErrInvalidUser
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{Username: username, PasswordHash: string(hash)}
	err = s.repo.CreateUser(ctx, &user)
	return user, err
}

// Храним в базе хеш токена. Сам токен возвращаем только при входе.
func tokenHash(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	if len(username) > 50 || len(password) > 72 {
		return "", ErrUnauthorized
	}
	user, err := s.repo.GetByUsername(ctx, strings.TrimSpace(username))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrUnauthorized
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)
	token := Token{Key: tokenHash(key), UserID: user.ID, ExpiresAt: time.Now().Add(24 * time.Hour)}
	if err := s.repo.CreateToken(ctx, &token); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Service) CurrentUser(ctx context.Context, key string) (User, error) {
	if len(key) != 64 {
		return User{}, ErrUnauthorized
	}
	if _, err := hex.DecodeString(key); err != nil {
		return User{}, ErrUnauthorized
	}
	user, err := s.repo.GetByToken(ctx, tokenHash(key))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUnauthorized
	}
	return user, err
}

func (s *Service) Logout(ctx context.Context, key string) error {
	return s.repo.DeleteToken(ctx, tokenHash(key))
}

func (s *Service) ChangeCredentials(ctx context.Context, key, currentPassword, username, newPassword string) (User, error) {
	user, err := s.CurrentUser(ctx, key)
	if err != nil {
		return User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)) != nil {
		return User{}, ErrUnauthorized
	}
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 50 || len(newPassword) < 8 || len(newPassword) > 72 {
		return User{}, ErrInvalidUser
	}
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			return User{}, ErrInvalidUser
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	if err := s.repo.UpdateCredentials(ctx, user.ID, username, string(hash)); err != nil {
		return User{}, err
	}
	user.Username = username
	user.PasswordHash = ""
	return user, nil
}
