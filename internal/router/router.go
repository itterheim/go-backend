package router

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repositories"
	"backend/internal/services"
	"backend/pkg/handler"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, config *config.AuthConfig) *http.ServeMux {
	r := http.NewServeMux()

	routes := make([]handler.Route, 0)

	tokenRepo := repositories.NewTokenRepository(db)
	userRepo := repositories.NewUserRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)

	// auth
	authService := auth.NewAuthService(userRepo, deviceRepo, tokenRepo, config.JWTSecret)
	authMiddleware := auth.GetAuthMiddleware(authService)
	var authHandler handler.Handler = handlers.NewAuthHandler(authService, config)
	routes = append(routes, authHandler.GetRoutes()...)

	// device
	deviceService := services.NewDeviceService(deviceRepo, authService)
	var deviceHandler handler.Handler = handlers.NewDeviceHandler(deviceService)
	routes = append(routes, deviceHandler.GetRoutes()...)

	// envi
	enviRepo := repositories.NewEnviRepository(db)
	var enviHandler handler.Handler = handlers.NewEnviHandler(enviRepo)
	routes = append(routes, enviHandler.GetRoutes()...)

	for _, route := range routes {
		if route.Public {
			r.Handle(route.Pattern, route.HandlerFunc)
		} else {
			r.Handle(route.Pattern, authMiddleware(route.HandlerFunc))
		}
	}
	// r.Handle("/envi/", http.StripPrefix("/envi", enviHandler.CreateRouter()))

	return r
}
