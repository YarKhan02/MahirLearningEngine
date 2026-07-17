package course

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	courseSvc  *Service
}

func NewHandler(courseSvc *Service) *Handler {
	return &Handler{courseSvc: courseSvc}
}

func (h *Handler) InsertCourse(c *gin.Context) {
	
	var req InsertCourse
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	course := ToCourse(req)

	createdCourse, err := h.courseSvc.InsertCourse(c.Request.Context(), course)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, ToCourseResponse(*createdCourse))
}

func (h *Handler) GetCourse(c *gin.Context) {
	
	courses, err := h.courseSvc.GetCourse(c.Request.Context())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]CourseResponse, 0, len(courses))
	for _, course := range courses {
		resp = append(resp, ToCourseResponse(course))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) DeleteCourse(c *gin.Context) {

	courseIDU, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id")
		return
	}

	err = h.courseSvc.DeleteCourse(c.Request.Context(), courseIDU)
	if err != nil {
		if errors.Is(err, ErrCourseNotFound) {
			response.WriteError(c, http.StatusNotFound, "course not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully deleted course")
}

func (h *Handler) InsertLesson(c *gin.Context) {

	courseID := c.Param("courseId")
	
	var req InsertLesson
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	courseIDU, err := uuid.Parse(courseID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid")
		return
	}

	req.CourseID = courseIDU

	lesson := ToLesson(req)

	err = h.courseSvc.InsertLesson(c.Request.Context(), lesson)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(c, http.StatusCreated, "lesson insert successfully")
}

func (h *Handler) GetLesson(c *gin.Context) {
	
	courseID := c.Param("courseId")

	courseIDU, err := uuid.Parse(courseID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid")
		return
	}

	lessons, err := h.courseSvc.GetLesson(c.Request.Context(), courseIDU)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]LessonResponse, 0, len(lessons))
	
	for _, lesson := range lessons {
		resp = append(resp, ToLessonResponse(lesson))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) UpdateLesson(c *gin.Context) {
	courseID := c.Param("courseId")
	lessonID := c.Param("lessonId")

	var req UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	courseIDU, err := uuid.Parse(courseID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid")
		return
	}
	req.CourseID = courseIDU

	lessonIDU, err := uuid.Parse(lessonID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid")
		return
	}
	req.ID = lessonIDU

	lesson := ToUpdateLesson(req)

	err = h.courseSvc.UpdateLesson(c.Request.Context(), lesson)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully updated")
}

func (h *Handler) ReorderLesson(c *gin.Context) {

	lessonID := c.Param("lessonId")
	
	var req UpdateLessonOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	
	lessonIDU, err := uuid.Parse(lessonID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid")
		return
	}

	err = h.courseSvc.ReorderLesson(c.Request.Context(), lessonIDU, req.OrderNo)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully reordered")
}