package liveclass

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc            *Service
	hub            *Hub
	originPatterns []string
}

func NewHandler(svc *Service, hub *Hub, originPatterns []string) *Handler {
	return &Handler{svc: svc, hub: hub, originPatterns: originPatterns}
}

/* Session lifecycle (admin) */

func (h *Handler) StartSession(c *gin.Context) {
	
	hostID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	
	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	
	batchID, err := uuid.Parse(req.BatchID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}
	
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id")
		return
	}

	sess, err := h.svc.StartSession(c.Request.Context(), hostID, batchID, courseID, strings.TrimSpace(req.Title))
	if err != nil {
		if errors.Is(err, ErrAlreadyLive) {
			response.WriteError(c, http.StatusConflict, "a class is already live for this batch")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	
	response.WriteJSON(c, http.StatusCreated, toSessionResponse(sess))
}

func (h *Handler) EndSession(c *gin.Context) {
	
	hostID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	
	sess, err := h.svc.EndSession(c.Request.Context(), id, hostID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "no live session to end")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	
	response.WriteJSON(c, http.StatusOK, toSessionResponse(sess))
}

/* Reads */

func (h *Handler) GetSession(c *gin.Context) {
	userID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	sess, err := h.svc.GetSessionForUser(c.Request.Context(), userID, role, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "session not found")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "no access to this class")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusOK, toSessionResponse(sess))
}

// GetMyLive returns the class currently live for the caller's own batch, if any.
func (h *Handler) GetMyLive(c *gin.Context) {
	userID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	sess, err := h.svc.GetLiveForUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "no live class")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toSessionResponse(sess))
}

// GetLiveByBatch lets a student find the class currently live for their batch.
func (h *Handler) GetLiveByBatch(c *gin.Context) {
	userID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}
	sess, err := h.svc.GetLiveByBatchForUser(c.Request.Context(), userID, role, batchID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "no live class for this batch")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "no access to this batch")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusOK, toSessionResponse(sess))
}

/* Student browse (past classes) */

func (h *Handler) MyBatches(c *gin.Context) {
	userID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.svc.MyBatches(c.Request.Context(), userID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toBatchOptions(items))
}

func (h *Handler) MyBatchCourses(c *gin.Context) {
	userID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}
	items, err := h.svc.MyBatchCourses(c.Request.Context(), userID, batchID)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toCourseOptions(items))
}

func (h *Handler) ListSessions(c *gin.Context) {
	userID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	batchID, err := uuid.Parse(c.Param("batchId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid batch id")
		return
	}
	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid course id")
		return
	}
	items, err := h.svc.SessionsForBatchCourse(c.Request.Context(), userID, role, batchID, courseID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.WriteError(c, http.StatusForbidden, "no access to this batch")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toSessionResponses(items))
}

/* WebSocket ticket + upgrade */

func (h *Handler) IssueTicket(c *gin.Context) {
	userID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	ticket, err := h.svc.IssueTicket(c.Request.Context(), userID, role, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "session not found")
		case errors.Is(err, ErrNotLive):
			response.WriteError(c, http.StatusConflict, "class is not live")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "no access to this class")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusOK, TicketResponse{Ticket: ticket})
}

// ServeWS is a public route (no auth middleware) — authorization is the
// one-time ticket minted by IssueTicket.
func (h *Handler) ServeWS(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}

	data, err := h.svc.ConsumeTicket(c.Request.Context(), c.Query("ticket"))
	
	if err != nil || data.SessionID != id {
		response.WriteError(c, http.StatusUnauthorized, "invalid or expired ticket")
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: h.originPatterns,
	})
	
	if err != nil {
		return // Accept already wrote the response
	}
	
	// Blocks until the client disconnects. Background ctx so it isn't cancelled
	// when the gin handler frame returns after the hijack.
	h.hub.Serve(context.Background(), conn, data.SessionID, data.UserID, data.Name, data.IsHost)
}

/* Video (LiveKit) */

func (h *Handler) RTCToken(c *gin.Context) {
	userID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	url, token, err := h.svc.RTCToken(c.Request.Context(), userID, role, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotConfigured):
			response.WriteError(c, http.StatusServiceUnavailable, "video is not enabled")
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "session not found")
		case errors.Is(err, ErrNotLive):
			response.WriteError(c, http.StatusConflict, "class is not live")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "no access to this class")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusOK, RTCTokenResponse{URL: url, Token: token})
}

func (h *Handler) SetMic(c *gin.Context) {
	hostID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	var req MicRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Identity == "" {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.svc.SetStudentMic(c.Request.Context(), hostID, role, id, req.Identity, req.Allow); err != nil {
		switch {
		case errors.Is(err, ErrVideoNotConfigured):
			response.WriteError(c, http.StatusServiceUnavailable, "video is not enabled")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "only the host can manage mics")
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "session not found")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusOK, "updated")
}

/* Snapshots */

func (h *Handler) SaveSnapshot(c *gin.Context) {
	userID, _, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	var req SaveSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	w, err := h.svc.SaveSnapshot(c.Request.Context(), id, userID, SaveSnapshot(req))
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalid):
			response.WriteError(c, http.StatusBadRequest, "empty whiteboard")
		case errors.Is(err, ErrForbidden):
			response.WriteError(c, http.StatusForbidden, "only the host can save")
		case errors.Is(err, ErrNotFound):
			response.WriteError(c, http.StatusNotFound, "session not found")
		default:
			response.WriteInternal(c, err)
		}
		return
	}
	response.WriteJSON(c, http.StatusCreated, toSnapshotResponse(w))
}

func (h *Handler) ListSnapshots(c *gin.Context) {
	userID, role, ok := middleware.CurrentUserRole(c)
	if !ok {
		response.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid session id")
		return
	}
	allowed, err := h.svc.CanView(c.Request.Context(), userID, role, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.WriteError(c, http.StatusNotFound, "session not found")
			return
		}
		response.WriteInternal(c, err)
		return
	}
	if !allowed {
		response.WriteError(c, http.StatusForbidden, "no access to this class")
		return
	}
	items, err := h.svc.ListSnapshots(c.Request.Context(), id)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}
	response.WriteJSON(c, http.StatusOK, toSnapshotResponses(items))
}
