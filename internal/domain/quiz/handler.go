package quiz

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

func (h *Handler) CreateQuiz(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.CreateQuiz(c.Request.Context(), lessonID, req); err != nil {
		switch {
		case errors.Is(err, ErrLessonNotFound):
			response.WriteError(c, http.StatusNotFound, "lesson not found")
		case errors.Is(err, ErrInvalidQuiz):
			response.WriteError(c, http.StatusBadRequest, "quiz must have a title and valid questions")
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusCreated, "quiz created")
}

func (h *Handler) ListQuizzes(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	quizzes, counts, err := h.svc.ListForLesson(c.Request.Context(), lessonID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	resp := make([]AdminQuizResponse, 0, len(quizzes))
	for _, q := range quizzes {
		cnt := counts[q.ID]
		resp = append(resp, ToAdminQuizResponse(q, cnt.Total, cnt.Pending))
	}
	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) EditQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.EditQuiz(c.Request.Context(), quizID, req); err != nil {
		switch {
		case errors.Is(err, ErrQuizNotFound):
			response.WriteError(c, http.StatusNotFound, "quiz not found")
		case errors.Is(err, ErrInvalidQuiz):
			response.WriteError(c, http.StatusBadRequest, "quiz must have a title and valid questions")
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusOK, "quiz updated")
}

func (h *Handler) DeleteQuiz(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	if err := h.svc.DeleteQuiz(c.Request.Context(), quizID); err != nil {
		if errors.Is(err, ErrQuizNotFound) {
			response.WriteError(c, http.StatusNotFound, "quiz not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "quiz deleted")
}

func (h *Handler) ListSubmissions(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	rows, summary, err := h.svc.SubmissionsForQuiz(c.Request.Context(), quizID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	items := make([]SubmissionRowResponse, 0, len(rows))
	for _, r := range rows {
		items = append(items, ToSubmissionRowResponse(r))
	}

	response.WriteJSON(c, http.StatusOK, gin.H{
		"submissions": items,
		"summary":     ToSubmissionSummaryResponse(summary),
	})
}

func (h *Handler) GetSubmission(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submissionId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid submission id")
		return
	}

	quiz, sub, name, err := h.svc.SubmissionDetail(c.Request.Context(), submissionID)
	if err != nil {
		if errors.Is(err, ErrSubmissionNotFound) {
			response.WriteError(c, http.StatusNotFound, "submission not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, ToSubmissionDetailResponse(quiz, sub, name))
}

func (h *Handler) GradeSubmission(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submissionId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid submission id")
		return
	}

	var req GradeQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	in, err := ToGradeInput(req, submissionID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Grade(c.Request.Context(), in); err != nil {
		if errors.Is(err, ErrSubmissionNotFound) {
			response.WriteError(c, http.StatusNotFound, "submission not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	response.WriteJSON(c, http.StatusOK, "submission graded")
}

func (h *Handler) ListMyQuizzes(c *gin.Context) {
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

	quizzes, err := h.svc.ListForStudent(c.Request.Context(), userID, lessonID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.WriteError(c, http.StatusForbidden, "no access to this lesson")
			return
		}
		response.WriteInternal(c, err)
		return
	}

	resp := make([]StudentQuizResponse, 0, len(quizzes))
	for _, q := range quizzes {
		resp = append(resp, ToStudentQuizResponse(q))
	}
	response.WriteJSON(c, http.StatusOK, resp)
}

func (h *Handler) SubmitQuiz(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	quizID, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	var req SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.svc.Submit(c.Request.Context(), userID, quizID, req); err != nil {
		switch {
		case errors.Is(err, ErrQuizNotFound):
			response.WriteError(c, http.StatusNotFound, "quiz not found")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "no access to this quiz")
		case errors.Is(err, ErrAlreadySubmitted):
			response.WriteError(c, http.StatusConflict, "you have already submitted this quiz")
		case errors.Is(err, ErrInvalidSubmission):
			response.WriteError(c, http.StatusBadRequest, "invalid answers")
		case errors.Is(err, ErrStudentNotFound):
			response.WriteError(c, http.StatusNotFound, "student profile not found")
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusOK, "quiz submitted")
}
