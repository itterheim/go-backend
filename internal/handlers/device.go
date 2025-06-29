package handlers

import (
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
		handler.NewRoute("GET /device/{id}/{$}", h.GetDevice, false, nil),
		handler.NewRoute("POST /device/{id}/{$}", h.CreateToken, false, nil),
	}
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.ListDevices()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
	}
	h.SendJSON(w, http.StatusOK, devices)
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
