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
		handler.NewRoute("GET /api/core/tags/{$}", h.ListTags, handler.RouteOwnerRole),
		handler.NewRoute("GET /api/core/tags/{tag}", h.GetTag, handler.RouteOwnerRole),
		handler.NewRoute("POST /api/core/tags", h.CreateTag, handler.RouteOwnerRole),
		handler.NewRoute("PUT /api/core/tags/{tag}", h.UpdateTag, handler.RouteOwnerRole),
		handler.NewRoute("DELETE /api/core/tags/{tag}", h.DeleteTag, handler.RouteOwnerRole),
	}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	private := r.URL.Query().Has("private")

	tags, err := h.service.ListTags(private)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusOK, tags)
}

func (h *TagHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetStringFromPath(r, "tag")
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
	var body CreateTagRequest
	err := h.ParseJSON(r, &body)
	if err != nil {
		h.SendJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.CreateTag(&body)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusCreated, result)
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	tag, err := h.GetStringFromPath(r, "tag")
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

	body.Tag = tag

	result, err := h.service.UpdateTag(&body)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusAccepted, result)
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tag, err := h.GetStringFromPath(r, "tag")
	if err != nil {
		h.SendJSON(w, http.StatusForbidden, err.Error())
		return
	}

	err = h.service.DeleteTag(tag)
	if err != nil {
		h.SendJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SendJSON(w, http.StatusAccepted, nil)
}
