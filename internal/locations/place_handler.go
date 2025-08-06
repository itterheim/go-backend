package locations

import (
	"backend/pkg/handler"
	"net/http"
)

type PlaceHandler struct {
	handler.BaseHandler

	service *PlaceService
}

func NewPlaceHandler(service *PlaceService) *PlaceHandler {
	return &PlaceHandler{service: service}
}

func (h *PlaceHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /api/locations/places/{$}", h.ListPlaces, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/locations/places/{id}", h.GetPlace, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/locations/places", h.CreatePlace, handler.RouteProviderRole),
		handler.NewRoute("PUT /api/locations/places/{id}", h.UpdatePlace, handler.RouteProviderRole),
		handler.NewRoute("DELETE /api/locations/places/{id}", h.DeletePlace, handler.RouteProviderRole),
	}
}

func (h *PlaceHandler) ListPlaces(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListPlaces()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *PlaceHandler) GetPlace(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.GetPlace(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if data == nil {
		h.SendJSON(w, http.StatusNotFound, "place not found")
		return
	}

	h.SendJSON(w, http.StatusOK, data)
}

func (h *PlaceHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	var data CreatePlaceRequest
	err := h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.CreatePlace(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *PlaceHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
	}

	var data Place
	err = h.ParseJSON(r, &data)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	data.ID = id

	result, err := h.service.UpdateHistory(&data)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, result)
}

func (h *PlaceHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
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
