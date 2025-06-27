package auth

import (
	"backend/pkg/auth"
	"backend/pkg/handler"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetAuthMiddleware(authService *AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := authenticate(r, authService)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), handler.RequestClaims, claims)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func authenticate(r *http.Request, authService *AuthService) (auth.Claims, error) {
	// GET /auth is public
	if r.URL.Path == "/auth" && r.Method == http.MethodPost {
		return auth.Claims{}, nil
	}

	// GET /auth/refresh is public, validation is done id the handler
	if r.URL.Path == "/auth/refresh" && r.Method == http.MethodGet {
		return auth.Claims{}, nil
	}

	claims, err := authenticateWithCookie(r, authService)
	if err != nil {
		claims, err = authenticateWithBearer(r, authService)
		return claims, err
	}

	return claims, err
}

func authenticateWithCookie(r *http.Request, authService *AuthService) (auth.Claims, error) {
	cookie, err := r.Cookie("access")
	if err != nil {
		return auth.Claims{}, err
	}

	// Check if the cookie is valid
	claims, err := authService.ValidateToken(cookie.Value)
	if err != nil {
		return auth.Claims{}, err
	}

	return claims, nil
}

func authenticateWithBearer(r *http.Request, authService *AuthService) (auth.Claims, error) {
	header := r.Header.Get("Authorization")
	if len(header) == 0 {
		return auth.Claims{}, errors.New("Missing Authorization header")
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return auth.Claims{}, errors.New("Invalid token")
	}

	token := strings.TrimPrefix(header, "Bearer ")

	claims, err := authService.ValidateToken(token)
	if err != nil {
		return claims, err
	}

	return claims, nil
}

// func authenticateWithBasicAuth(r *http.Request) bool {
// 	username, password, ok := r.BasicAuth()
// 	if !ok {
// 		return false
// 	}

// 	fmt.Println(username, password)

// 	return true
// }
