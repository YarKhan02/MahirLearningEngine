package student

import (
	"errors"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/common"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/crypto"
	"github.com/YarKhan02/MahirLearningEngine/internal/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	studentSvc   *Service
	userSvc      common.StudentAccountRegistrar
	tempPassword string
}

func NewHandler(studentSvc *Service, userSvc common.StudentAccountRegistrar, tempPassword string) *Handler {
	return &Handler{
		studentSvc:   studentSvc,
		userSvc:      userSvc,
		tempPassword: tempPassword,
	}
}

func (h *Handler) RegisterStudent(c *gin.Context) {

	var req RegisterStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	batchIDU, err := uuid.Parse(req.BatchID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	s, err := ToRegisterStudent(req)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid date of birth")
		return
	}

	err = h.studentSvc.RegisterStudent(c.Request.Context(), s, batchIDU)
	if err != nil {
		if errors.Is(err, ErrUsernameAlreadyRegistered) {
			response.WriteError(c, http.StatusConflict, "this username is already taken")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, "successfully registered")
}

func (h *Handler) GetStudents(c *gin.Context) {

	p := pagination.Parse(c.Query("page"), c.Query("pageSize"), 10, 10)

	students, total, err := h.studentSvc.GetStudents(c.Request.Context(), c.Query("q"), p.Limit(), p.Offset())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	items := make([]AdminStudentResponse, 0, len(students))
	for _, s := range students {
		items = append(items, ToAdminStudentResponse(s))
	}

	response.WriteJSON(c, http.StatusOK, pagination.NewPage(items, total, p))
}

func (h *Handler) UpdateStudentStatus(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req UpdateStudentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.studentSvc.UpdateStudentStatus(c.Request.Context(), studentIDU, req.Status); err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully updated status")
}

func (h *Handler) UpdateStudentBatch(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req UpdateStudentBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	var batchID *uuid.UUID
	if req.BatchID != "" {
		batchIDU, err := uuid.Parse(req.BatchID)
		if err != nil {
			response.WriteError(c, http.StatusBadRequest, "invalid batch id")
			return
		}
		batchID = &batchIDU
	}

	if err := h.studentSvc.UpdateStudentBatch(c.Request.Context(), studentIDU, batchID); err != nil {
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "successfully updated batch")
}

func (h *Handler) CreateStudentAccount(c *gin.Context) {

	studentIDU, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid student id")
		return
	}

	s, err := h.studentSvc.GetStudentByID(c.Request.Context(), studentIDU)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			response.WriteError(c, http.StatusNotFound, "student not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	password, err := crypto.GenerateTempPassword(h.tempPassword, 10)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	if err := h.userSvc.RegisterStudentAccount(c.Request.Context(), s.Username, password); err != nil {
		if errors.Is(err, common.ErrUsernameTaken) {
			response.WriteError(c, http.StatusConflict, "this student already has an account")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, StudentAccountResponse{
		Username: s.Username,
		Password: password,
	})
}

func (h *Handler) AdminCreateStudent(c *gin.Context) {

	var req RegisterStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	batchIDU, err := uuid.Parse(req.BatchID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}

	s, err := ToRegisterStudent(req)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid date of birth")
		return
	}

	err = h.studentSvc.RegisterStudent(c.Request.Context(), s, batchIDU)
	if err != nil {
		if errors.Is(err, ErrUsernameAlreadyRegistered) {
			response.WriteError(c, http.StatusConflict, "this username is already taken")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusCreated, "successfully created student")
}

func (h *Handler) GetMyCourses(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courses, err := h.studentSvc.GetStudentCourses(c.Request.Context(), userID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]StudentCourseResponse, 0, len(courses))
	for _, course := range courses {
		resp = append(resp, ToStudentCourseResponse(course))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) GetMyLessons(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	courseIDU, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id")
		return
	}

	lessons, err := h.studentSvc.GetStudentLessons(c.Request.Context(), userID, courseIDU)
	if err != nil {
		if errors.Is(err, ErrCourseAccessDenied) {
			response.WriteError(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, ErrStudentNotFound) {
			response.WriteError(c, http.StatusNotFound, "student profile not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	resp := make([]StudentLessonResponse, 0, len(lessons))
	for _, lesson := range lessons {
		resp = append(resp, ToStudentLessonResponse(lesson))
	}

	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) SetLessonProgress(c *gin.Context) {

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonIDU, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	var req SetLessonProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.studentSvc.SetLessonProgress(c.Request.Context(), userID, lessonIDU, *req.Completed); err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			response.WriteError(c, http.StatusNotFound, "student profile not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "progress updated")
}
