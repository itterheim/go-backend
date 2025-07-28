package core

import (
	"backend/pkg/handler"
	"net/http"
)

type UserHandler struct {
	handler.BaseHandler

	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /api/core/users/{$}", h.ListUsers, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/core/users/{id}", h.GetUser, handler.RouteOwnerRole),
	}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		h.SendJSON(w, http.StatusNotFound, "User not found")
		return
	}

	h.SendJSON(w, http.StatusOK, user)
}
