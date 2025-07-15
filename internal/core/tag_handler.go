package core

import (
	"backend/pkg/handler"
	"net/http"
)

type TagHandler struct {
	handler.BaseHandler

	service *TagService
}

func NewTagHandler(service *TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /core/tags/{$}", h.ListTags, handler.RouteAuthenticatedRole),
	}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.ListTags()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, tags)
}
