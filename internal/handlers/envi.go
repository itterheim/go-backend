package handlers

import (
	"backend/internal/repositories"
	"backend/pkg/handler"
	"net/http"
)

type EnviHandler struct {
	handler.BaseHandler

	repo *repositories.Envi
}

func NewEnviHandler(repo *repositories.Envi) *EnviHandler {
	return &EnviHandler{repo: repo}
}

func (h *EnviHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /envi/{$}", h.GetEnvi, false, nil),
		handler.NewRoute("GET /envi/public/{$}", h.GetEnviPublic, true, nil),
	}
}

func (h *EnviHandler) GetEnvi(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.SendJSON(w, http.StatusOK, claims.ID)
}

func (h *EnviHandler) GetEnviPublic(w http.ResponseWriter, r *http.Request) {
	h.SendJSON(w, http.StatusOK, "/envi/public")
}
