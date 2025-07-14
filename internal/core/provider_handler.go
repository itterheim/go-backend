package core

import (
	"backend/pkg/handler"
	"encoding/json"
	"net/http"
	"time"
)

type ProviderHandler struct {
	handler.BaseHandler

	service *ProviderService
}

func NewProviderHandler(service *ProviderService) *ProviderHandler {
	return &ProviderHandler{
		service: service,
	}
}

func (h *ProviderHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /core/providers/{$}", h.GetProviders, handler.RouteOwnerRole),
		handler.NewRoute("POST /core/providers/{$}", h.CreateProvider, handler.RouteOwnerRole),
		handler.NewRoute("GET /core/providers/{id}/{$}", h.GetProvider, handler.RouteOwnerRole),
		handler.NewRoute("PUT /core/providers/{id}/{$}", h.UpdateProvider, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /core/providers/{id}/{$}", h.DeleteProvider, handler.RouteOwnerRole),

		handler.NewRoute("POST /core/providers/{id}/token/{$}", h.CreateToken, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /core/providers/{id}/token/{$}", h.RevokeToken, handler.RouteOwnerRole),
	}
}

func (h *ProviderHandler) GetProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.service.ListProviders()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
	}
	h.SendJSON(w, http.StatusOK, providers)
}

func (h *ProviderHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	body := struct {
		name        string
		description string
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "invalid provider data")
		return
	}

	providers, err := h.service.CreateProvider(body.name, body.description)
	if err != nil {
		// TODO: differentiate between validation errors and general server errors
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusCreated, providers)
}

func (h *ProviderHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	provider, err := h.service.GetProvider(providerId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	if provider == nil {
		h.SendJSON(w, http.StatusNotFound, "provider not found")
		return
	}

	h.SendJSON(w, http.StatusOK, provider)
}

func (h *ProviderHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	body := Provider{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	provider, err := h.service.UpdateProvider(&body)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusOK, provider)
}

func (h *ProviderHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.DeleteProvider(providerId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *ProviderHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
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

	token, err := h.service.CreateToken(providerId, lifespan)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJSON(w, http.StatusCreated, token)
}

func (h *ProviderHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.RevokeToken(providerId)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
