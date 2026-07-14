package handler

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/mapper"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/crypto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentHandler struct {
	studentSvc 		*student.Service
	userSvc    		*user.Service
	tempPassword 	string
}

func NewStudentHandler(studentSvc *student.Service, userSvc *user.Service, tempPassword string) *StudentHandler {
	return &StudentHandler{
		studentSvc: 	studentSvc, 
		userSvc: 		userSvc,
		tempPassword: 	tempPassword,
	}
}

func (h *StudentHandler) RegisterStudent(c *gin.Context) {

	var req dto.RegisterStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	batchIDU, err := uuid.Parse(req.BatchID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	s, err := mapper.ToRegisterStudent(req)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid date of birth")
		return
	}

	err = h.studentSvc.RegisterStudent(c.Request.Context(), s, batchIDU)
	if err != nil {
		if errors.Is(err, student.ErrUsernameAlreadyRegistered) {
			writeError(c, http.StatusConflict, "this username is already taken")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, "successfully registered")
}

func (h *StudentHandler) GetStudents(c *gin.Context) {

	students, err := h.studentSvc.GetStudents(c.Request.Context(), c.Query("q"))
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.AdminStudentResponse, 0, len(students))
	for _, s := range students {
		resp = append(resp, mapper.ToAdminStudentResponse(s))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *StudentHandler) UpdateStudentStatus(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req dto.UpdateStudentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.studentSvc.UpdateStudentStatus(c.Request.Context(), studentIDU, req.Status); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully updated status")
}

func (h *StudentHandler) UpdateStudentBatch(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req dto.UpdateStudentBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	var batchID *uuid.UUID
	if req.BatchID != "" {
		batchIDU, err := uuid.Parse(req.BatchID)
		if err != nil {
			writeError(c, http.StatusBadRequest, "invalid batch id")
			return
		}
		batchID = &batchIDU
	}

	if err := h.studentSvc.UpdateStudentBatch(c.Request.Context(), studentIDU, batchID); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "successfully updated batch")
}

func (h *StudentHandler) CreateStudentAccount(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	s, err := h.studentSvc.GetStudentByID(c.Request.Context(), studentIDU)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			writeError(c, http.StatusNotFound, "student not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	password, err := crypto.GenerateTempPassword(h.tempPassword, 10)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to generate password")
		return
	}

	if _, err := h.userSvc.RegisterStudentAccount(c.Request.Context(), s.Username, password); err != nil {
		if errors.Is(err, user.ErrUsernameTaken) {
			writeError(c, http.StatusConflict, "this student already has an account")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, dto.StudentAccountResponse{
		Username: s.Username,
		Password: password,
	})
}

func (h *StudentHandler) AdminCreateStudent(c *gin.Context) {

	var req dto.RegisterStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	batchIDU, err := uuid.Parse(req.BatchID)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	s, err := mapper.ToRegisterStudent(req)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid date of birth")
		return
	}

	err = h.studentSvc.RegisterStudent(c.Request.Context(), s, batchIDU)
	if err != nil {
		if errors.Is(err, student.ErrUsernameAlreadyRegistered) {
			writeError(c, http.StatusConflict, "this username is already taken")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, "successfully created student")
}

func (h *StudentHandler) GetMyCourses(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courses, err := h.studentSvc.GetStudentCourses(c.Request.Context(), userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.StudentCourseResponse, 0, len(courses))
	for _, course := range courses {
		resp = append(resp, mapper.ToStudentCourseResponse(course))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *StudentHandler) GetMyLessons(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courseIDU, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid course id")
		return
	}

	lessons, err := h.studentSvc.GetStudentLessons(c.Request.Context(), userID, courseIDU)
	if err != nil {
		if errors.Is(err, student.ErrCourseAccessDenied) {
			writeError(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, student.ErrStudentNotFound) {
			writeError(c, http.StatusNotFound, "student profile not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.StudentLessonResponse, 0, len(lessons))
	for _, lesson := range lessons {
		resp = append(resp, mapper.ToStudentLessonResponse(lesson))
	}

	writeJSON(c, http.StatusOK, resp)
}

func (h *StudentHandler) SetLessonProgress(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	var req dto.SetLessonProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.studentSvc.SetLessonProgress(c.Request.Context(), userID, lessonIDU, *req.Completed); err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			writeError(c, http.StatusNotFound, "student profile not found")
			return
		}
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, "progress updated")
}
