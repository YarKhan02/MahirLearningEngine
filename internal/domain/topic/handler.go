package topic

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

func (h *Handler) CreateTopic(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	var req InsertTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.CreateTopic(c.Request.Context(), lessonID, req); err != nil {
		if errors.Is(err, ErrLessonNotFound) {
			response.WriteError(c, http.StatusNotFound, "lesson not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, "topic created")
}

func (h *Handler) ListTopics(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	topics, err := h.svc.GetTopics(c.Request.Context(), lessonID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, toResponses(topics))
}

func (h *Handler) ListMyTopics(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	topics, err := h.svc.GetTopicsForStudent(c.Request.Context(), userID, lessonID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.WriteError(c, http.StatusForbidden, "no access to this lesson")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, toResponses(topics))
}

func (h *Handler) UpdateTopic(c *gin.Context) {
	topicID, err := uuid.Parse(c.Param("topicId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid topic id")
		return
	}

	var req UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.UpdateTopic(c.Request.Context(), ToUpdateTopic(req, topicID)); err != nil {
		if errors.Is(err, ErrTopicNotFound) {
			response.WriteError(c, http.StatusNotFound, "topic not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "topic updated")
}

func (h *Handler) ReorderTopic(c *gin.Context) {
	topicID, err := uuid.Parse(c.Param("topicId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid topic id")
		return
	}

	var req UpdateTopicOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.ReorderTopic(c.Request.Context(), topicID, req.OrderNo); err != nil {
		switch {
		case errors.Is(err, ErrTopicNotFound):
			response.WriteError(c, http.StatusNotFound, "topic not found")
		case errors.Is(err, ErrInvalidOrderNo):
			response.WriteError(c, http.StatusBadRequest, "order out of range")
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusOK, "topic reordered")
}

func (h *Handler) DeleteTopic(c *gin.Context) {
	topicID, err := uuid.Parse(c.Param("topicId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid topic id")
		return
	}

	if err := h.svc.DeleteTopic(c.Request.Context(), topicID); err != nil {
		if errors.Is(err, ErrTopicNotFound) {
			response.WriteError(c, http.StatusNotFound, "topic not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "topic deleted")
}

func toResponses(topics []Topic) []TopicResponse {
	resp := make([]TopicResponse, 0, len(topics))
	for _, t := range topics {
		resp = append(resp, ToTopicResponse(t))
	}
	return resp
}
