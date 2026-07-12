package handler

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/batch"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BatchHandler struct {
	batchSvc *batch.Service
}

func NewBatchHandler(batchSvc *batch.Service) *BatchHandler {
	return &BatchHandler{batchSvc: batchSvc}
}

func (h *BatchHandler) CreateBatch(c *gin.Context) {

	var req dto.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	batch, err := mapper.ToCreateBatch(req)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "invalid")
		return
	}

	err = h.batchSvc.CreateBatch(c.Request.Context(), batch)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, "succesfully created batch")
}

func (h *BatchHandler) GetBatches(c *gin.Context) {

	batches, err := h.batchSvc.GetBatches(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.BatchResponse, 0, len(batches))
	for _, batch := range batches {
		resp = append(resp, mapper.ToBatchResponse(batch))
	}

	writeJSON(c, http.StatusOK, resp)
}
func (h *BatchHandler) GetBatchCourses(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	courses, err := h.batchSvc.GetBatchCourses(c.Request.Context(), batchIDU)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.BatchCourseResponse, 0, len(courses))
	for _, course := range courses {
		resp = append(resp, mapper.ToBatchCourseResponse(course))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *BatchHandler) UpdateBatchCourses(c *gin.Context) {

	batchIDU, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	var req dto.UpdateBatchCoursesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	add, err := helper.ParseUUIDs(req.AddCourseIDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid course id in addCourseIds")
		return
	}

	remove, err := helper.ParseUUIDs(req.RemoveCourseIDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid course id in removeCourseIds")
		return
	}

	var grantedBy *uuid.UUID
	if value, exists := c.Get("claims"); exists {
		if claims, ok := value.(*token.Claims); ok {
			if userID, err := uuid.Parse(claims.UserID); err == nil {
				grantedBy = &userID
			}
		}
	}

	err = h.batchSvc.UpdateBatchCourses(c.Request.Context(), batchIDU, add, remove, grantedBy)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully updated batch courses")
}

func (h *BatchHandler) GetPublicBatches(c *gin.Context) {

	batches, err := h.batchSvc.GetOpenBatchesWithCourses(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.PublicBatchResponse, 0, len(batches))
	for _, b := range batches {
		resp = append(resp, mapper.ToPublicBatchResponse(b))
	}

	writeJSON(c, http.StatusOK, resp)
}
