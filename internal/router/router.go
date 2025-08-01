package router

import (
	"backend/internal/config"
	"backend/internal/core"
	"backend/internal/locations"
	"backend/pkg/handler"
	"net/http"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, config *config.AuthConfig) *http.ServeMux {
	r := http.NewServeMux()

	routes := make([]handler.Route, 0)

	// repositories
	tokenRepo := core.NewTokenRepository(db)
	userRepo := core.NewUserRepository(db)
	providerRepo := core.NewProviderRepository(db)
	eventRepo := core.NewEventRepository(db)
	tagRepo := core.NewTagRepository(db)

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

	// tags
	tagService := core.NewTagService(tagRepo)
	var tagHandler handler.Handler = core.NewTagHandler(tagService)
	routes = append(routes, tagHandler.GetRoutes()...)

	// location - history
	locationRepo := locations.NewLocationRepository(db)
	locationService := locations.NewLocationService(locationRepo, eventRepo)
	var locationHandler handler.Handler = locations.NewLocationHandler(locationService)
	routes = append(routes, locationHandler.GetRoutes()...)

	// location - places
	placeRepo := locations.NewPlaceRepository(db)
	placeService := locations.NewPlaceService(placeRepo)
	var placeHandler handler.Handler = locations.NewPlaceHandler(placeService)
	routes = append(routes, placeHandler.GetRoutes()...)

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

	fs := http.FileServer(http.Dir("./web"))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := http.Dir("./web").Open(r.URL.Path)
		if err == nil {
			defer file.Close()
			fileInfo, err := file.Stat()
			if err == nil && !fileInfo.IsDir() {
				fs.ServeHTTP(w, r)
				return
			}
		}

		indexPath := filepath.Join("web", "index.html")
		http.ServeFile(w, r, indexPath)
	})

	return r
}
