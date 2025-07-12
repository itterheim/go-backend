package handler

import "net/http"

type RouteRole string

const (
	RouteOwnerRole         = "Owner"         // Only the main user can use the route
	RouteProviderRole      = "Provider"      // Main user and provider can use this route
	RouteAuthenticatedRole = "Authenticated" // anyone/anything with auth token
	RoutePublicRole        = "Public"        // public route without authentication
)

type Route struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
	Role        RouteRole
}

func NewRoute(pattern string, handlerFunc http.HandlerFunc, role RouteRole) Route {
	return Route{
		Pattern:     pattern,
		HandlerFunc: handlerFunc,
		Role:        role,
	}
}
