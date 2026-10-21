package auth

import (
	"errors"
	"net/http"
	"strings"

	"laboratory_work_3/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type credentialsRequest struct {
	CurrentPassword string `json:"current_password"`
	Username        string `json:"username"`
	NewPassword     string `json:"new_password"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/users", h.Register)
	mux.HandleFunc("POST /api/auth/token/login", h.Login)
	mux.Handle("POST /api/auth/token/logout", h.RequireAuth(http.HandlerFunc(h.Logout)))
	mux.Handle("GET /api/auth/users/me", h.RequireAuth(http.HandlerFunc(h.CurrentUser)))
	mux.Handle("PATCH /api/auth/users/me", h.RequireAuth(http.HandlerFunc(h.ChangeCredentials)))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidUser)
		return
	}
	w.Header().Set("Location", "/api/auth/users/me")
	httpapi.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	key, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, map[string]string{"auth_token": key})
}

func requestToken(r *http.Request) string {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) == 2 && (strings.EqualFold(parts[0], "Token") || strings.EqualFold(parts[0], "Bearer")) {
		return parts[1]
	}
	return ""
}

func handleAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrUnauthorized) {
		w.Header().Set("WWW-Authenticate", "Token")
		httpapi.WriteError(w, http.StatusUnauthorized, "требуется действующий токен или верные данные для входа")
		return
	}
	httpapi.HandleError(w, err, ErrInvalidUser)
}

// Каждый запрос к учебным данным должен содержать токен.
func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.service.CurrentUser(r.Context(), requestToken(r)); err != nil {
			handleAuthError(w, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.CurrentUser(r.Context(), requestToken(r))
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Logout(r.Context(), requestToken(r)); err != nil {
		handleAuthError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChangeCredentials(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.ChangeCredentials(r.Context(), requestToken(r), req.CurrentPassword, req.Username, req.NewPassword)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, user)
}
