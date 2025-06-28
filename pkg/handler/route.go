package handler

import "net/http"

type Route struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
	Permissions *[]string
	Public      bool
}

func NewRoute(pattern string, handlerFunc http.HandlerFunc, public bool, permissions *[]string) Route {
	return Route{
		Pattern:     pattern,
		HandlerFunc: handlerFunc,
		Permissions: permissions,
		Public:      public,
	}
}
