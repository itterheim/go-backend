package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/pkg/handler"
	"encoding/json"
	"net/http"
	"time"
)

type DeviceHandler struct {
	handler.BaseHandler

	service *services.Device
}

func NewDeviceHandler(service *services.Device) *DeviceHandler {
	return &DeviceHandler{
		service: service,
	}
}

func (h *DeviceHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /device/{$}", h.GetDevices, false, nil),
		handler.NewRoute("POST /device/{$}", h.CreateDevice, false, nil),

		handler.NewRoute("GET /device/{id}/{$}", h.GetDevice, false, nil),
		handler.NewRoute("PUT /device/{id}/{$}", h.UpdateDevice, false, nil),
		handler.NewRoute("DELETE /device/{id}/{$}", h.DeleteDevice, false, nil),

		handler.NewRoute("POST /device/{id}/token/{$}", h.CreateToken, false, nil),
		handler.NewRoute("DELETE /device/{id}/token/{$}", h.RevokeToken, false, nil),
	}
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.ListDevices()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
	}
	h.SendJSON(w, http.StatusOK, devices)
}

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	body := struct {
		name        string
		description string
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "invalid device data")
		return
	}

	device, err := h.service.CreateDevice(body.name, body.description)
	if err != nil {
		// TODO: differentiate between validation errors and general server errors
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusCreated, device)
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	device, err := h.service.GetDevice(deviceId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	if device == nil {
		h.SendJSON(w, http.StatusNotFound, "device not found")
		return
	}

	h.SendJSON(w, http.StatusOK, device)
}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	body := models.Device{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	device, err := h.service.UpdateDevice(&body)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, device)
}

func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.DeleteDevice(deviceId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *DeviceHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	deviceId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	// lifespan in hours
	body := struct{ lifespan int }{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "Invalid request body: lifespan")
		return
	}

	lifespan := time.Hour * 24 * time.Duration(body.lifespan)

	token, err := h.service.CreateToken(deviceId, lifespan)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusCreated, token)
}

func (h *DeviceHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	deviceId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.RevokeToken(deviceId)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
