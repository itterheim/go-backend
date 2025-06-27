package handler

import (
	"backend/pkg/auth"
	"encoding/json"
	"fmt"
	"net/http"
)

const RequestClaims = "claims"

type BaseHandler struct {
}

// Retrieve claims from the request context
func (h *BaseHandler) GetClaimsFromContext(r *http.Request) (auth.Claims, error) {
	value, ok := r.Context().Value(RequestClaims).(auth.Claims)
	if !ok {
		return auth.Claims{}, fmt.Errorf("no claims found in context")
	}

	return value, nil
}

// respondWithJSON writes the given data as JSON response with the specified status code.
func (h *BaseHandler) SendJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// parseJSON decodes the JSON request body into the given interface.
func (h *BaseHandler) ParseJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}