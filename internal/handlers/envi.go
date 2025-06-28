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

func (h *EnviHandler) CreateRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /{$}", h.GetEnvi)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}

func (h *EnviHandler) GetEnvi(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.SendJSON(w, http.StatusOK, claims.ID)
}
