package liveclass

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSession(ctx context.Context, s LiveSession) error
	EndSession(ctx context.Context, id, hostID uuid.UUID) error
	GetSession(ctx context.Context, id uuid.UUID) (LiveSession, error)
	GetLiveByBatch(ctx context.Context, batchID uuid.UUID) (LiveSession, error)
	GetLiveForUser(ctx context.Context, userID uuid.UUID) (LiveSession, error)
	StudentCanJoin(ctx context.Context, userID, sessionID uuid.UUID) (bool, error)

	CreateSnapshot(ctx context.Context, w WhiteboardSnapshot) error
	NextSnapshotVersion(ctx context.Context, sessionID uuid.UUID) (int, error)
	ListSnapshots(ctx context.Context, sessionID uuid.UUID) ([]WhiteboardSnapshot, error)

	// Student browse (past classes by batch → course).
	SessionsByBatchCourse(ctx context.Context, batchID, courseID uuid.UUID) ([]LiveSession, error)
	StudentBatches(ctx context.Context, userID uuid.UUID) ([]BatchOption, error)
	StudentBatchCourses(ctx context.Context, userID, batchID uuid.UUID) ([]CourseOption, error)
	StudentInBatch(ctx context.Context, userID, batchID uuid.UUID) (bool, error)

	UserDisplayName(ctx context.Context, userID uuid.UUID) (string, error)
}
