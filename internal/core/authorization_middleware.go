package core

import (
	"backend/pkg/handler"
	"backend/pkg/jwt"
	"fmt"
	"net/http"
)

func GetAuthorizationMiddleware(authService *AuthService) func(http.Handler, handler.RouteRole) http.Handler {
	return func(next http.Handler, requiredRole handler.RouteRole) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(handler.RequestClaims).(jwt.Claims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			fmt.Println("Authorize", claims)

			allowed := false

			switch requiredRole {
			case handler.RouteOwnerRole:
				if claims.Type == jwt.UserClaim && claims.Role == jwt.OwnerRole {
					allowed = true
				}
			case handler.RouteProviderRole:
				if claims.Type == jwt.ProviderClaim || (claims.Type == jwt.UserClaim && claims.Role == jwt.OwnerRole) {
					allowed = true
				}
			}

			if allowed {
				// Call the next handler
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		})
	}
}
