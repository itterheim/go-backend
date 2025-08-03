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
		handler.NewRoute("GET /api/core/providers/{$}", h.GetProviders, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/core/providers", h.CreateProvider, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/core/providers/{id}", h.GetProvider, handler.RouteOwnerRole),
		handler.NewRoute("PUT /api/core/providers/{id}", h.UpdateProvider, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /api/core/providers/{id}", h.DeleteProvider, handler.RouteOwnerRole),

		handler.NewRoute("POST /api/core/providers/{id}/token", h.CreateToken, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /api/core/providers/{id}/token", h.RevokeToken, handler.RouteOwnerRole),
	}
}

func (h *ProviderHandler) GetProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.service.ListProviders()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.SendJSON(w, http.StatusOK, providers)
}

func (h *ProviderHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "invalid provider data\n"+err.Error())
		return
	}

	provider, err := h.service.CreateProvider(body.Name, body.Description)
	if err != nil {
		// TODO: differentiate between validation errors and general server errors
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, provider)
}

func (h *ProviderHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	provider, err := h.service.GetProvider(providerId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if provider == nil {
		h.SendJSON(w, http.StatusNotFound, "provider not found")
		return
	}

	h.SendJSON(w, http.StatusOK, provider)
}

func (h *ProviderHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	provider, err := h.service.UpdateProvider(providerId, body.Name, body.Description)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, provider)
}

func (h *ProviderHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteProvider(providerId)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *ProviderHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	// lifespan in days
	var body struct {
		Lifespan int `json:"lifespan"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, "Invalid request body: lifespan")
		return
	}

	lifespan := time.Hour * 24 * time.Duration(body.Lifespan)

	token, err := h.service.CreateToken(claims.UserID, providerId, lifespan)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, token)
}

func (h *ProviderHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	providerId, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.RevokeToken(providerId)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
