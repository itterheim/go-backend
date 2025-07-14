package locations

import (
	"backend/pkg/handler"
	"net/http"
	"time"
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
		handler.NewRoute("GET /locations/history/{$}", h.ListHistory, handler.RouteOwnerRole),
		handler.NewRoute("GET /locations/history/{id}", h.GetHistory, handler.RouteOwnerRole),
		handler.NewRoute("POST /locations/history", h.RegisterHistory, handler.RouteProviderRole),
		handler.NewRoute("PUT /locations/history/{id}", h.UpdateHistory, handler.RouteProviderRole),
		handler.NewRoute("DELETE /locations/history/{id}", h.DeleteHistory, handler.RouteProviderRole),
	}
}

func (h *LocationHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	// TODO: read from query string
	from := time.UnixMilli(0)
	to := time.Now()

	data, err := h.service.ListHistory(from, to)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
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

	var data CreateGpsHistoryRequest
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.UserID = claims.UserID
	data.ProviderID = claims.ProviderID

	result, err := h.service.RegisterHistory(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *LocationHandler) UpdateHistory(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
	}

	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	var data GpsHistory
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ID = id

	result, err := h.service.UpdateHistory(claims.UserID, &data)
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
