package maintenance

import (
	"context"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"go.uber.org/zap"
)

type Service struct {
	repo  Repository
	store ObjectStore
}

func NewService(repo Repository, store ObjectStore) *Service {
	return &Service{repo: repo, store: store}
}

func (s *Service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *Service) RunAll(ctx context.Context) Report {
	var rep Report
	log := logging.FromLogger(ctx)

	if n, err := s.repo.PurgeExpiredTokens(ctx); err != nil {
		log.Warn("purge expired tokens failed", zap.String("event", "maintenance_tokens_failed"), zap.Error(err))
	} else {
		rep.ExpiredTokensDeleted = n
	}

	stale, err := s.repo.ListStalePendingAttachments(ctx)
	if err != nil {
		log.Warn("list stale attachments failed", zap.String("event", "maintenance_attachments_failed"), zap.Error(err))
		return rep
	}
	for _, a := range stale {
		// Delete the R2 object first. If that fails, keep the DB row so the next
		// run retries — never orphan the object by dropping the row prematurely.
		if err := s.store.DeleteObject(ctx, a.R2Key); err != nil {
			rep.R2DeleteFailures++
			continue
		}
		if err := s.repo.DeleteAttachmentByID(ctx, a.ID); err != nil {
			continue
		}
		rep.PendingAttachmentsSwept++
	}

	log.Info("maintenance cleanup complete",
		zap.String("event", "maintenance_cleanup"),
		zap.Int64("tokens_deleted", rep.ExpiredTokensDeleted),
		zap.Int("attachments_swept", rep.PendingAttachmentsSwept),
		zap.Int("r2_failures", rep.R2DeleteFailures),
	)
	return rep
}
