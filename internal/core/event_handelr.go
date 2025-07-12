package core

import (
	"backend/pkg/handler"
	"net/http"
	"time"
)

type EventHandler struct {
	handler.BaseHandler

	service *EventService
}

func NewEventHandler(service *EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /core/events/{$}", h.GetEvents, handler.RouteOwnerRole),
		handler.NewRoute("GET /core/events/{id}", h.GetEvent, handler.RouteOwnerRole),
		handler.NewRoute("PUT /core/events/{id}", h.UpdateEvent, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /core/events/{id}", h.DeleteEvent, handler.RouteOwnerRole),
		handler.NewRoute("POST /core/events", h.CreateEvent, handler.RouteProviderRole),
	}
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: read from query string
	from := time.UnixMilli(0)
	to := time.Now()

	data, err := h.service.ListEvents(from, to)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}
	h.SendJSON(w, http.StatusOK, data)
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	event, err := h.service.GetEvent(eventId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	if event == nil {
		h.SendJSON(w, http.StatusNotFound, "event not found")
		return
	}

	h.SendJSON(w, http.StatusOK, event)
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: both user and provider in claims
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err)
		return
	}

	var data CreateEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	data.UserID = claims.UserID

	result, err := h.service.CreateEvent(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, result)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	var data UpdateEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	data.ID = eventId

	result, err := h.service.UpdateEvent(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
	}

	h.SendJSON(w, http.StatusOK, result)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.DeleteEvent(eventId)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusOK)
}
