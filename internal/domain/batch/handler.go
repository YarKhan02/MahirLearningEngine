package batch

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/YarKhan02/MahirLearningEngine/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateBatch(c *gin.Context) {

	var req CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	batch, err := ToCreateBatch(req)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid")
		return
	}

	err = h.svc.CreateBatch(c.Request.Context(), batch)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, "successfully created batch")
}

func (h *Handler) UpdateBatch(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req UpdateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	updated, err := ToUpdateBatch(batchIDU, req)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.svc.UpdateBatch(c.Request.Context(), updated)
	if err != nil {
		if errors.Is(err, ErrBatchNotFound) {
			response.WriteError(c, http.StatusNotFound, "batch not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully updated batch")
}

func (h *Handler) DeleteBatch(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	err = h.svc.DeleteBatch(c.Request.Context(), batchIDU)
	if err != nil {
		if errors.Is(err, ErrBatchNotFound) {
			response.WriteError(c, http.StatusNotFound, "batch not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully deleted batch")
}

func (h *Handler) GetBatches(c *gin.Context) {

	batches, err := h.svc.GetBatches(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]BatchResponse, 0, len(batches))
	for _, batch := range batches {
		resp = append(resp, ToBatchResponse(batch))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}
func (h *Handler) GetBatchCourses(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	courses, err := h.svc.GetBatchCourses(c.Request.Context(), batchIDU)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]BatchCourseResponse, 0, len(courses))
	for _, course := range courses {
		resp = append(resp, ToBatchCourseResponse(course))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) UpdateBatchCourses(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req UpdateBatchCoursesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	add, err := helper.ParseUUIDs(req.AddCourseIDs)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id in addCourseIds")
		return
	}

	remove, err := helper.ParseUUIDs(req.RemoveCourseIDs)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id in removeCourseIds")
		return
	}

	var grantedBy *uuid.UUID
	if userID, ok := middleware.CurrentUserID(c); ok {
		grantedBy = &userID
	}

	err = h.svc.UpdateBatchCourses(c.Request.Context(), batchIDU, add, remove, grantedBy)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully updated batch courses")
}

func (h *Handler) GetPublicBatches(c *gin.Context) {

	batches, err := h.svc.GetOpenBatchesWithCourses(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]PublicBatchResponse, 0, len(batches))
	for _, b := range batches {
		resp = append(resp, ToPublicBatchResponse(b))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}
