package liveclass

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/r2"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	roleAdmin      = "admin"
	ticketTTL      = 30 * time.Second
	ticketKeyPart  = "wsticket:"
)

type Service struct {
	repo    Repository
	r2      *r2.Client
	redis   *redis.RedisClient
	livekit *LiveKit
}

func NewService(repo Repository, r2 *r2.Client, redis *redis.RedisClient, livekit *LiveKit) *Service {
	return &Service{repo: repo, r2: r2, redis: redis, livekit: livekit}
}

// TicketData is what a consumed WS ticket resolves to.
type TicketData struct {
	UserID    uuid.UUID `json:"userId"`
	SessionID uuid.UUID `json:"sessionId"`
	Name      string    `json:"name"`
	IsHost    bool      `json:"isHost"`
}

// Session lifecycle 

func (s *Service) StartSession(ctx context.Context, hostID, batchID, courseID uuid.UUID, title string) (LiveSession, error) {
	// One live session per batch. The partial unique index is the real guard;
	// this check gives a clean error in the common case.
	if _, err := s.repo.GetLiveByBatch(ctx, batchID); err == nil {
		return LiveSession{}, ErrAlreadyLive
	} else if !errors.Is(err, ErrNotFound) {
		return LiveSession{}, err
	}

	sess := LiveSession{
		ID:       uuid.New(),
		BatchID:  batchID,
		CourseID: courseID,
		HostID:   hostID,
		Title:    title,
		Status:   StatusLive,
	}
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return LiveSession{}, err
	}

	full, err := s.repo.GetSession(ctx, sess.ID)
	if err != nil {
		return LiveSession{}, err
	}
	logging.FromLogger(ctx).Info("live session started",
		zap.String("event", "live_session_started"),
		zap.String("session_id", full.ID.String()),
		zap.String("batch_id", batchID.String()),
		zap.String("host_id", hostID.String()),
	)
	return full, nil
}

func (s *Service) EndSession(ctx context.Context, id, hostID uuid.UUID) (LiveSession, error) {
	if err := s.repo.EndSession(ctx, id, hostID); err != nil {
		return LiveSession{}, err
	}
	logging.FromLogger(ctx).Info("live session ended",
		zap.String("event", "live_session_ended"),
		zap.String("session_id", id.String()),
	)
	return s.repo.GetSession(ctx, id)
}

func (s *Service) GetSession(ctx context.Context, id uuid.UUID) (LiveSession, error) {
	return s.repo.GetSession(ctx, id)
}

func (s *Service) GetLiveByBatch(ctx context.Context, batchID uuid.UUID) (LiveSession, error) {
	return s.repo.GetLiveByBatch(ctx, batchID)
}

// GetLiveForUser resolves the live class for the caller's own batch (student).
func (s *Service) GetLiveForUser(ctx context.Context, userID uuid.UUID) (LiveSession, error) {
	return s.repo.GetLiveForUser(ctx, userID)
}

// Authorization 

// AuthorizeJoin reports whether the user may join the session and whether they
// are the host (the one allowed to draw/clear/save). Returns ErrNotFound,
// ErrNotLive, or ErrForbidden.
func (s *Service) AuthorizeJoin(ctx context.Context, userID uuid.UUID, role string, sessionID uuid.UUID) (bool, error) {
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return false, err
	}
	if sess.Status != StatusLive {
		return false, ErrNotLive
	}
	if role == roleAdmin {
		return sess.HostID == userID, nil
	}
	ok, err := s.repo.StudentCanJoin(ctx, userID, sessionID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, ErrForbidden
	}
	return false, nil
}

// CanView reports whether the user may view a session's saved boards, ignoring
// live status (students revisit boards after class). Returns ErrNotFound.
func (s *Service) CanView(ctx context.Context, userID uuid.UUID, role string, sessionID uuid.UUID) (bool, error) {
	if role == roleAdmin {
		if _, err := s.repo.GetSession(ctx, sessionID); err != nil {
			return false, err
		}
		return true, nil
	}
	return s.repo.StudentCanJoin(ctx, userID, sessionID)
}

// WebSocket tickets 

// IssueTicket authorizes the user and mints a one-time, short-lived ticket the
// client uses to open the WebSocket (browsers can't set auth headers on WS).
func (s *Service) IssueTicket(ctx context.Context, userID uuid.UUID, role string, sessionID uuid.UUID) (string, error) {
	isHost, err := s.AuthorizeJoin(ctx, userID, role, sessionID)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	ticket := hex.EncodeToString(buf)

	name, _ := s.repo.UserDisplayName(ctx, userID) // best-effort; empty is fine

	data, err := json.Marshal(TicketData{UserID: userID, SessionID: sessionID, Name: name, IsHost: isHost})
	if err != nil {
		return "", err
	}
	if err := s.redis.Set(ctx, ticketKeyPart+ticket, string(data), ticketTTL); err != nil {
		return "", err
	}
	return ticket, nil
}

// ConsumeTicket validates and immediately invalidates a ticket (single use).
func (s *Service) ConsumeTicket(ctx context.Context, ticket string) (TicketData, error) {
	if ticket == "" {
		return TicketData{}, ErrForbidden
	}
	key := ticketKeyPart + ticket
	raw, err := s.redis.Get(ctx, key)
	if err != nil || raw == "" {
		return TicketData{}, ErrForbidden
	}
	_ = s.redis.Delete(ctx, key) // one-time use

	var data TicketData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return TicketData{}, ErrForbidden
	}
	return data, nil
}

// Video (LiveKit)

// RTCToken authorizes the caller for the session and mints a LiveKit join token.
// The host may publish (camera/mic/screen); students subscribe only until the
// teacher lets them speak. Returns ErrVideoNotConfigured when LiveKit is off.
func (s *Service) RTCToken(ctx context.Context, userID uuid.UUID, role string, sessionID uuid.UUID) (url string, token string, err error) {
	isHost, err := s.AuthorizeJoin(ctx, userID, role, sessionID)
	if err != nil {
		return "", "", err
	}
	name, _ := s.repo.UserDisplayName(ctx, userID)
	token, err = s.livekit.Token(userID.String(), name, sessionID.String(), isHost)
	if err != nil {
		return "", "", err
	}
	return s.livekit.WSURL(), token, nil
}

// SetStudentMic lets the host grant/revoke a student's ability to speak.
func (s *Service) SetStudentMic(ctx context.Context, hostID uuid.UUID, role string, sessionID uuid.UUID, studentIdentity string, allow bool) error {
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if role != roleAdmin || sess.HostID != hostID {
		return ErrForbidden
	}
	return s.livekit.SetPublish(ctx, sessionID.String(), studentIdentity, allow)
}

// Snapshots

func (s *Service) SaveSnapshot(ctx context.Context, sessionID, userID uuid.UUID, snap SaveSnapshot) (WhiteboardSnapshot, error) {
	if len(snap.Scene) == 0 {
		return WhiteboardSnapshot{}, ErrInvalid
	}
	// Only the host may save.
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return WhiteboardSnapshot{}, err
	}
	if sess.HostID != userID {
		return WhiteboardSnapshot{}, ErrForbidden
	}

	version, err := s.repo.NextSnapshotVersion(ctx, sessionID)
	if err != nil {
		return WhiteboardSnapshot{}, err
	}

	key := fmt.Sprintf("live/%s/whiteboards/%d.excalidraw", sessionID.String(), version)
	if err := s.r2.PutObject(ctx, key, "application/json", bytes.NewReader(snap.Scene)); err != nil {
		return WhiteboardSnapshot{}, err
	}

	w := WhiteboardSnapshot{
		ID:         uuid.New(),
		SessionID:  sessionID,
		StorageKey: key,
		ImageURL:   snap.ImageURL,
		Version:    version,
		CreatedBy:  userID,
	}
	if err := s.repo.CreateSnapshot(ctx, w); err != nil {
		return WhiteboardSnapshot{}, err
	}

	logging.FromLogger(ctx).Info("whiteboard saved",
		zap.String("event", "whiteboard_saved"),
		zap.String("session_id", sessionID.String()),
		zap.Int("version", version),
	)
	return w, nil
}

func (s *Service) ListSnapshots(ctx context.Context, sessionID uuid.UUID) ([]WhiteboardSnapshot, error) {
	return s.repo.ListSnapshots(ctx, sessionID)
}

/* ---------------- Student browse (past classes) ---------------- */

func (s *Service) MyBatches(ctx context.Context, userID uuid.UUID) ([]BatchOption, error) {
	return s.repo.StudentBatches(ctx, userID)
}

func (s *Service) MyBatchCourses(ctx context.Context, userID, batchID uuid.UUID) ([]CourseOption, error) {
	return s.repo.StudentBatchCourses(ctx, userID, batchID)
}

// SessionsForBatchCourse lists a batch+course's classes (newest first). Admins
// see any; students only their own batches.
func (s *Service) SessionsForBatchCourse(ctx context.Context, userID uuid.UUID, role string, batchID, courseID uuid.UUID) ([]LiveSession, error) {
	if role != roleAdmin {
		ok, err := s.repo.StudentInBatch(ctx, userID, batchID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrForbidden
		}
	}
	return s.repo.SessionsByBatchCourse(ctx, batchID, courseID)
}
