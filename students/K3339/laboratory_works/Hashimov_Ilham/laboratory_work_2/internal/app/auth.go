package app

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"sync"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]int64
}

type Auth struct {
	userRepo UserRepository
	sessions *SessionManager
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]int64),
	}
}

func NewAuth(userRepo UserRepository, sessions *SessionManager) *Auth {
	return &Auth{
		userRepo: userRepo,
		sessions: sessions,
	}
}

func (s *SessionManager) Create(userID int64) (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	sessionID := hex.EncodeToString(bytes)

	s.mu.Lock()
	s.sessions[sessionID] = userID
	s.mu.Unlock()

	return sessionID, nil
}

func (s *SessionManager) Get(sessionID string) (int64, bool) {
	s.mu.RLock()
	userID, ok := s.sessions[sessionID]
	s.mu.RUnlock()

	return userID, ok
}

func (s *SessionManager) Delete(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/login.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
			return
		}

		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := a.userRepo.GetByUsername(r.Context(), username)
		if err != nil {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		if user.PasswordHash != password {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		sessionID, err := a.sessions.Create(user.ID)
		if err != nil {
			http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if user.IsAdmin {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/homeworks", http.StatusSeeOther)
		return
	}
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		a.sessions.Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *Auth) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, ok := a.sessions.Get(cookie.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := a.userRepo.GetByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
			return
		}

		if !user.IsAdmin {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func (a *Auth) CurrentUser(r *http.Request) (User, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return User{}, false
	}

	userID, ok := a.sessions.Get(cookie.Value)
	if !ok {
		return User{}, false
	}

	user, err := a.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		return User{}, false
	}

	return user, true
}
