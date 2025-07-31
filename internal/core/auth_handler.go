package core

import (
	"backend/internal/config"
	"backend/pkg/handler"
	"net/http"
)

type AuthHandler struct {
	*handler.BaseHandler

	config  *config.AuthConfig
	service *AuthService
}

func NewAuthHandler(service *AuthService, config *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		config:  config,
		service: service,
	}
}

func (h *AuthHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /api/auth", h.Validate, handler.RouteProviderRole),
		handler.NewRoute("POST /api/auth", h.Login, handler.RoutePublicRole),
		handler.NewRoute("DELETE /api/auth", h.Logout, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/auth/refresh", h.Refresh, handler.RoutePublicRole),
	}
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	responseData := map[string]any{
		"userId":     claims.UserID,
		"providerId": claims.ProviderID,
		"type":       claims.Type,
	}

	h.SendJSON(w, http.StatusOK, responseData)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Get the username and password from the request body
	type LoginRequestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequestData
	err := h.ParseJSON(r, &loginRequest)
	if err != nil {
		h.removeCookies(w)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, accessToken, err := h.service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		h.removeCookies(w)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.setCookies(w, refreshToken, accessToken)

	// Send a response
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.removeCookies(w)

	// Send a response
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from the cookie
	cookie, err := r.Cookie("refresh")
	if err != nil {
		http.Error(w, "no refresh token found", http.StatusUnauthorized)
		return
	}

	refreshToken, accessToken, err := h.service.ValidateRefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.removeCookies(w)
	h.setCookies(w, refreshToken, accessToken)

	// Send a response
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) setCookies(w http.ResponseWriter, refreshToken, accessToken string) {
	refreshCookie := &http.Cookie{
		Name:     "refresh",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   h.config.Secure,
		SameSite: http.SameSiteNoneMode,
		Path:     "/api/auth/refresh",
	}
	http.SetCookie(w, refreshCookie)

	// Create an authentication cookie
	accessCookie := &http.Cookie{
		Name:     "access",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   h.config.Secure,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, accessCookie)
}

func (h *AuthHandler) removeCookies(w http.ResponseWriter) {
	expiredCookie := &http.Cookie{
		Name:     "refresh",
		Value:    "",
		HttpOnly: true,
		Secure:   h.config.Secure,
		Path:     "/api/auth/refresh",
		MaxAge:   -1,
	}
	http.SetCookie(w, expiredCookie)

	expiredCookie = &http.Cookie{
		Name:     "access",
		Value:    "",
		HttpOnly: true,
		Secure:   h.config.Secure,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(w, expiredCookie)
}
