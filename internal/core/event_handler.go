package core

import (
	"backend/pkg/handler"
	"net/http"
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
		handler.NewRoute("GET /api/core/events/{$}", h.GetEvents, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/core/events/{id}", h.GetEvent, handler.RouteOwnerRole),
		handler.NewRoute("PUT /api/core/events/{id}", h.UpdateEvent, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /api/core/events/{id}", h.DeleteEvent, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/core/events", h.CreateEvent, handler.RouteProviderRole),
	}
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	query := &EventQueryBuilder{}
	err := query.FromRequest(r)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
	}

	data, err := h.service.ListEvents(query)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.SendJSON(w, http.StatusOK, data)
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.service.GetEvent(eventId)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if event == nil {
		h.SendJSON(w, http.StatusNotFound, "event not found")
		return
	}

	h.SendJSON(w, http.StatusOK, event)
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	var data CreateEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ProviderID = claims.ProviderID

	result, err := h.service.CreateEvent(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	var data UpdateEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ID = eventId
	data.ProviderID = claims.ProviderID

	result, err := h.service.UpdateEvent(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, result)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteEvent(eventId)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusOK)
}
