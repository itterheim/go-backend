package raw

import (
	"backend/internal/core"
	"backend/pkg/handler"
	"net/http"
)

type RawHandler struct {
	handler.BaseHandler
	service *RawService
}

func NewRawHandler(service *RawService) *RawHandler {
	return &RawHandler{service: service}
}

func (h *RawHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /api/raw/{$}", h.ListRawEvents, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/raw/{id}", h.GetRawEvent, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/raw", h.CreateRawEvent, handler.RouteProviderRole),
		handler.NewRoute("PUT /api/raw/{id}", h.UpdateRawEvent, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /api/raw/{id}", h.DeleteRawEvent, handler.RouteOwnerRole),
	}
}

func (h *RawHandler) ListRawEvents(w http.ResponseWriter, r *http.Request) {
	query := &core.EventQueryBuilder{}
	err := query.FromRequest(r)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.ListRawEvents(query)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *RawHandler) GetRawEvent(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.GetRawEvent(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if data == nil {
		h.SendJSON(w, http.StatusNotFound, "raw data not found")
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *RawHandler) CreateRawEvent(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
	}

	var data CreateRawEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ProviderID = claims.ProviderID

	result, err := h.service.RegisterRawEvent(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *RawHandler) UpdateRawEvent(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
	}

	var data UpdateRawEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ID = id
	data.ProviderID = claims.ProviderID

	result, err := h.service.UpdateRawEvent(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, result)
}

func (h *RawHandler) DeleteRawEvent(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteRawEvent(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
