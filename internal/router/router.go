package router

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repositories"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, config *config.AuthConfig) *http.ServeMux {
	r := http.NewServeMux()

	tokenRepo := repositories.NewTokenRepository(db)
	userRepo := repositories.NewUserRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)

	authService := auth.NewAuthService(userRepo, deviceRepo, tokenRepo, config.JWTSecret)
	authMiddleware := auth.GetAuthMiddleware(authService)
	authHandler := handlers.NewAuthHandler(authService, config)

	authHandler.RegisterRoutes(r, authMiddleware)
	// r.Handle("/auth/", http.StripPrefix("/auth", authHandler.CreateRouter(authMiddleware)))

	// var notesHandler handler.Handler = notes.NewNotesHandler()
	// r.Handle("/notes/", http.StripPrefix("/notes", authMiddleware(notesHandler.CreateRouter())))

	// enviRepo := repositories.NewEnviRepository(db)
	// var enviHandler handler.Handler = handlers.NewEnviHandler(enviRepo)
	// r.Handle("/envi/", http.StripPrefix("/envi", enviHandler.CreateRouter()))

	return r
}
