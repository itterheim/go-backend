package handler

import (
	"backend/pkg/jwt"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const RequestClaims = "claims"

type BaseHandler struct {
}

// Retrieve claims from the request context
func (h *BaseHandler) GetClaimsFromContext(r *http.Request) (jwt.Claims, error) {
	value, ok := r.Context().Value(RequestClaims).(jwt.Claims)
	if !ok {
		return jwt.Claims{}, fmt.Errorf("no claims found in context")
	}

	return value, nil
}

func (h *BaseHandler) GetUserClaimsFromContext(r *http.Request) (jwt.Claims, error) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		return jwt.Claims{}, err
	}

	if claims.Type == jwt.UserClaim {
		return claims, nil
	}

	return jwt.Claims{}, errors.New("jwt claims are not of UserClaim type")
}

func (h *BaseHandler) GetAuthIdsFromContext(r *http.Request) (int64, *int64, error) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		return 0, nil, err
	}

	return claims.UserID, claims.ProviderID, nil
}

// Parse in64 id form the request path
func (h *BaseHandler) GetInt64FromPath(r *http.Request, key string) (int64, error) {
	id := r.PathValue("id")
	if len(id) == 0 {
		return 0, errors.New("invalid provider id")
	}

	providerId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.New("invalid provider id")
	}

	return providerId, nil
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
