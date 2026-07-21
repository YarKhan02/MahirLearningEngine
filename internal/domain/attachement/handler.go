package attachement

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
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

func (h *Handler) PresignUpload(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req PresignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if !AllowedTypes[req.ContentType] {
		response.WriteError(c, http.StatusBadRequest, "unsupported file type")
		return
	}
	if req.SizeBytes <= 0 || req.SizeBytes > MaxUploadSize {
		response.WriteError(c, http.StatusBadRequest, "file too large")
		return
	}

	url, err := h.svc.PresignUpload(c.Request.Context(), userID, ToPresignRequestEntity(req))
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrFailed):
			response.WriteError(c, http.StatusBadRequest, err.Error())
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusOK, ToPresignResponseDTO(url))
}

func (h *Handler) ConfirmUpload(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Key == "" {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	a, err := h.svc.ConfirmUpload(c.Request.Context(), userID, req.Key)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "upload not found")
		case errors.Is(err, ErrFailed):
			response.WriteError(c, http.StatusBadRequest, "file too large")
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusOK, ToAttachmentResponse(*a))
}

func (h *Handler) ListCourseMaterials(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id")
		return
	}

	items, err := h.svc.ListCourseMaterials(c.Request.Context(), courseID.String())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, toResponses(items))
}

func (h *Handler) ListMyCourseMaterials(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id")
		return
	}

	items, err := h.svc.ListCourseMaterialsForStudent(c.Request.Context(), userID, courseID.String())
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.WriteError(c, http.StatusForbidden, "no access to this course")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, toResponses(items))
}

func (h *Handler) DeleteMaterial(c *gin.Context) {
	id, err := uuid.Parse(c.Param("attachmentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid attachment id")
		return
	}

	if err := h.svc.DeleteMaterial(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "attachment not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "material deleted")
}

func toResponses(items []Attachment) []AttachmentResponse {
	resp := make([]AttachmentResponse, 0, len(items))
	for _, a := range items {
		resp = append(resp, ToAttachmentResponse(a))
	}
	return resp
}
