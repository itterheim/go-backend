package core

import (
	"backend/pkg/handler"
	"net/http"
)

type ActionHandler struct {
	handler.BaseHandler

	actionService *ActionService
	eventService  *EventService
}

func NewActionHandler(action *ActionService, event *EventService) *ActionHandler {
	return &ActionHandler{actionService: action, eventService: event}
}

func (h *ActionHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /core/actions/{$}", h.GetActions, handler.RouteOwnerRole),
		handler.NewRoute("GET /core/action/{id}", h.GetAction, handler.RouteOwnerRole),
		handler.NewRoute("PUT /core/action/{id}", h.UpdateAction, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /core/action/{id}", h.DeleteAction, handler.RouteOwnerRole),
		handler.NewRoute("POST /core/action", h.CreateAction, handler.RouteProviderRole),
	}
}

func (h *ActionHandler) GetActions(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err)
		return
	}

	data, err := h.actionService.ListActions(claims.UserID)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *ActionHandler) GetAction(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err)
		return
	}

	actionId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	data, err := h.actionService.GetAction(actionId, claims.UserID)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *ActionHandler) UpdateAction(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err)
		return
	}

	actionId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	var data struct {
		action UpdateActionRequest
		event  *CreateEventRequest
	}
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, data)
	}

	if data.event != nil {
		data.event.UserID = claims.UserID
		event, err := h.eventService.CreateEvent(data.event)
		if err != nil {
			h.SendJSON(w, http.StatusInternalServerError, err)
			return
		}
		data.action.EventID = event.ID
	}

	data.action.ID = actionId

	response, err := h.actionService.UpdateAction(&data.action)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, response)
}

func (h *ActionHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err)
		return
	}

	actionId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.actionService.DeleteAction(actionId, claims.UserID)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ActionHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetUserClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err)
		return
	}

	var data struct {
		action CreateActionRequest
		event  *CreateEventRequest
	}
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	if data.event != nil {
		data.event.UserID = claims.UserID
		event, err := h.eventService.CreateEvent(data.event)
		if err != nil {
			h.SendJSON(w, http.StatusInternalServerError, err)
		}
		data.action.EventID = event.ID
	}

	result, err := h.actionService.CreateAction(&data.action)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, result)
}
