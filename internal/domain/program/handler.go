package program

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Public

func (h *Handler) ListPublished(c *gin.Context) {
	items, err := h.svc.ListPublished(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toResponses(items))
}

func (h *Handler) GetBySlug(c *gin.Context) {
	p, err := h.svc.GetPublishedBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "program not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toResponse(p))
}

// Admin

func (h *Handler) ListAll(c *gin.Context) {
	items, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toResponses(items))
}

func (h *Handler) Create(c *gin.Context) {
	var req UpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), toUpsert(req))
	if err != nil {
		if errors.Is(err, ErrInvalid) {
			response.WriteError(c, http.StatusUnprocessableEntity, "title is required")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusCreated, toResponse(p))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid program id")
		return
	}
	var req UpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, toUpsert(req))
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalid):
			response.WriteError(c, http.StatusUnprocessableEntity, "title is required")
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "program not found")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusOK, toResponse(p))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid program id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "program not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, "program deleted")
}

type setPublishedRequest struct {
	Published bool `json:"published"`
}

func (h *Handler) SetPublished(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid program id")
		return
	}
	var req setPublishedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	p, err := h.svc.SetPublished(c.Request.Context(), id, req.Published)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "program not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toResponse(p))
}

func (h *Handler) Reorder(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			response.WriteError(c, http.StatusBadRequest, "invalid program id in list")
			return
		}
		ids = append(ids, id)
	}
	if err := h.svc.Reorder(c.Request.Context(), ids); err != nil {
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, "reordered")
}
