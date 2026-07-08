package handler

import (
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/course"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	course := mapper.ToCourse(req)

	createdCourse, err := h.courseSvc.InsertCourse(c.Request.Context(), course)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "course insertion failed")
		return
	}

	writeJSON(c, http.StatusCreated, mapper.ToCourseResponse(*createdCourse))
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	
	courses, err := h.courseSvc.GetCourse(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "course fetch failed")
		return
	}

	resp := make([]dto.CourseResponse, 0, len(courses))
	for _, course := range courses {
		resp = append(resp, mapper.ToCourseResponse(course))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *CourseHandler) InsertLesson(c *gin.Context) {

	courseID := c.Param("courseId")
	
	var req dto.InsertLesson
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	courseIDU, err := uuid.Parse(courseID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid")
		return
	}

	req.CourseID = courseIDU

	lesson := mapper.ToLesson(req)

	err = h.courseSvc.InsertLesson(c.Request.Context(), lesson)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, "lesson insert successfully")
}

func (h *CourseHandler) GetLesson(c *gin.Context) {
	
	courseID := c.Param("courseId")

	courseIDU, err := uuid.Parse(courseID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid")
		return
	}

	lessons, err := h.courseSvc.GetLesson(c.Request.Context(), courseIDU)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "lesson fetch failed")
		return
	}

	resp := make([]dto.LessonResponse, 0, len(lessons))
	
	for _, lesson := range lessons {
		resp = append(resp, mapper.ToLessonResponse(lesson))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *CourseHandler) UpdateLesson(c *gin.Context) {
	courseID := c.Param("courseId")
	lessonID := c.Param("lessonId")

	var req dto.UpdateLesson
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	courseIDU, err := uuid.Parse(courseID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid")
		return
	}
	req.CourseID = courseIDU

	lessonIDU, err := uuid.Parse(lessonID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid")
		return
	}
	req.ID = lessonIDU

	lesson := mapper.ToUpdateLesson(req)

	err = h.courseSvc.UpdateLesson(c.Request.Context(), lesson)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "succesfully updated")
}

func (h *CourseHandler) ReorderLesson(c *gin.Context) {

	lessonID := c.Param("lessonId")
	
	var req dto.UpdateLessonOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	
	lessonIDU, err := uuid.Parse(lessonID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid")
		return
	}

	err = h.courseSvc.ReorderLesson(c.Request.Context(), lessonIDU, req.OrderNo)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully reordered")
}