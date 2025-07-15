package router

import (
	"backend/internal/config"
	"backend/internal/core"
	"backend/internal/locations"
	"backend/pkg/handler"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, config *config.AuthConfig) *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/", fs)

	routes := make([]handler.Route, 0)

	// repositories
	tokenRepo := core.NewTokenRepository(db)
	userRepo := core.NewUserRepository(db)
	providerRepo := core.NewProviderRepository(db)
	eventRepo := core.NewEventRepository(db)

	// auth
	authService := core.NewAuthService(userRepo, providerRepo, tokenRepo, config)
	authenticationMiddleware := core.GetAuthenticationMiddleware(authService)
	authorizationMiddleware := core.GetAuthorizationMiddleware(authService)
	var authHandler handler.Handler = core.NewAuthHandler(authService, config)
	routes = append(routes, authHandler.GetRoutes()...)

	// users
	userService := core.NewUserService(userRepo)
	var userHandler handler.Handler = core.NewUserHandler(userService)
	routes = append(routes, userHandler.GetRoutes()...)

	// provider
	providerService := core.NewProviderService(providerRepo, authService)
	var providerHandler handler.Handler = core.NewProviderHandler(providerService)
	routes = append(routes, providerHandler.GetRoutes()...)

	// events
	eventService := core.NewEventService(eventRepo)
	var eventHandler handler.Handler = core.NewEventHandler(eventService)
	routes = append(routes, eventHandler.GetRoutes()...)

	// location
	locationRepo := locations.NewLocationRepository(db)
	locationService := locations.NewLocationService(locationRepo, eventRepo)
	var locationHandler handler.Handler = locations.NewLocationHandler(locationService)
	routes = append(routes, locationHandler.GetRoutes()...)

	for _, route := range routes {
		var handlerFunc http.Handler = route.HandlerFunc

		if route.Role == handler.RouteOwnerRole || route.Role == handler.RouteProviderRole {
			handlerFunc = authorizationMiddleware(handlerFunc, route.Role)
		}

		if route.Role != handler.RoutePublicRole {
			handlerFunc = authenticationMiddleware(handlerFunc)
		}

		r.Handle(route.Pattern, handlerFunc)
	}

	return r
}
