package handler

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseSvc  *course.Service
}

func NewCourseHandler(courseSvc *course.Service) *CourseHandler {
	return &CourseHandler{courseSvc: courseSvc}
}

func (h *CourseHandler) InsertCourse(c *gin.Context) {
	
	var req dto.InsertCourse
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	course, err := h.courseSvc.InsertCourse(c.Request.Context(), req)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "course insertion failed")
		return
	}

	cr := dto.CourseResponse{
		ID: course.ID,
		Title: course.Title,
		Level: course.Level,
		Duration: course.Duration,
		Description: course.Description,
	}

	if course.IsActive {
		cr.Status = "active"
	} else {
		cr.Status = "archived"
	}

	writeJSON(c, http.StatusCreated, cr)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	courses, err := h.courseSvc.GetCourse(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "course fetch failed")
		return
	}

	resp := make([]dto.CourseResponse, 0, len(courses))
	
	for _, course := range courses {
		status := "archived"
		if course.IsActive {
			status = "active"
		}

		resp = append(resp, dto.CourseResponse{
			ID:          course.ID,
			Title:       course.Title,
			Level:       course.Level,
			Duration:    course.Duration,
			Description: course.Description,
			Status:      status,
		})
	}

	writeJSON(c, http.StatusOK, resp)
}