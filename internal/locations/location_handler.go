package locations

import (
	"backend/internal/core"
	"backend/pkg/handler"
	"net/http"
)

type LocationHandler struct {
	handler.BaseHandler

	service *LocationService
}

func NewLocationHandler(service *LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /api/locations/history/{$}", h.ListHistory, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/locations/history/{id}", h.GetHistory, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/locations/history", h.RegisterHistory, handler.RouteProviderRole),
		handler.NewRoute("PUT /api/locations/history/{id}", h.UpdateHistory, handler.RouteProviderRole),
		handler.NewRoute("DELETE /api/locations/history/{id}", h.DeleteHistory, handler.RouteProviderRole),
	}
}

func (h *LocationHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	query := &core.EventQueryBuilder{}
	err := query.FromRequest(r)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.ListHistory(query)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *LocationHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.GetHistory(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if data == nil {
		h.SendJSON(w, http.StatusNotFound, "gps history not found")
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *LocationHandler) RegisterHistory(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	var data CreateLocationEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ProviderID = claims.ProviderID

	result, err := h.service.RegisterHistory(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *LocationHandler) UpdateHistory(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
	}

	var data UpdateLocationEventRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ID = id
	data.ProviderID = claims.ProviderID

	result, err := h.service.UpdateHistory(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, result)
}

func (h *LocationHandler) DeleteHistory(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteHistory(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
