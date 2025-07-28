package core

import (
	"backend/pkg/handler"
	"net/http"
)

type TagHandler struct {
	handler.BaseHandler

	service *TagService
}

func NewTagHandler(service *TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) GetRoutes() []handler.Route {
	return []handler.Route{
		handler.NewRoute("GET /api/core/tags/{$}", h.ListTags, handler.RouteAuthenticatedRole),
		handler.NewRoute("GET /api/core/tags/{id}", h.GetTag, handler.RouteAuthenticatedRole),
		handler.NewRoute("POST /api/core/tags", h.CreateTag, handler.RouteOwnerRole),
		handler.NewRoute("PUT /api/core/tags/{id}", h.UpdateTag, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /api/core/tags/{id}", h.DeleteTag, handler.RouteOwnerRole),
	}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.ListTags()
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, tags)
}

func (h *TagHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
	}

	tag, err := h.service.GetTag(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tag == nil {
		h.SendJSON(w, http.StatusNotFound, "tag not found")
	}

	h.SendJSON(w, http.StatusOK, tag)
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	var body CreateTagRequest
	err = h.ParseJSON(r, &body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	body.UserID = claims.UserID

	result, err := h.service.CreateTag(&body)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	claims, err := h.GetClaimsFromContext(r)
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	var body UpdateTagRequest
	err = h.ParseJSON(r, &body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	body.ID = id
	body.UserID = claims.UserID

	result, err := h.service.UpdateTag(&body)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusAccepted, result)
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetInt64FromPath(r, "id")
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	err = h.service.DeleteTag(id)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusAccepted, nil)
}
