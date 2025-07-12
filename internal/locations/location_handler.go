package locations

import "backend/pkg/handler"

type LocationHandler struct {
	handler.BaseHandler

	service *LocationService
}

func NewLocationHandler(service *LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) GetRoutes() []handler.Route {
	return []handler.Route{}
}
