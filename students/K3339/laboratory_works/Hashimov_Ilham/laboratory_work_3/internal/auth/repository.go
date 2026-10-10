package auth

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *Repository) CreateToken(ctx context.Context, token *Token) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *Repository) GetByToken(ctx context.Context, key string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Joins("JOIN api_tokens ON api_tokens.user_id = api_users.id").
		Where("api_tokens.key = ? AND api_tokens.expires_at > ?", key, time.Now()).
		First(&user).Error
	return user, err
}

func (r *Repository) DeleteToken(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&Token{}).Error
}
