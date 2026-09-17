package codeexec

import (
	"context"

	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/redis"
	"github.com/google/uuid"
)

type Service struct {
	runner *Runner
	redis  *redis.RedisClient
}

func NewService(runner *Runner, redis *redis.RedisClient) *Service {
	return &Service{runner: runner, redis: redis}
}

// Run validates, rate-limits, and executes a single interactive run for a user.
func (s *Service) Run(ctx context.Context, userID uuid.UUID, source, stdin string) (RunResult, error) {
	if len(source) == 0 {
		return RunResult{}, ErrInvalid
	}
	
	if len(source) > maxSourceBytes || len(stdin) > maxStdinBytes {
		return RunResult{}, ErrTooLarge
	}
	
	if ok, err := s.redis.AcquireLock(ctx, "coderun:"+userID.String(), minRunInterval); err == nil && !ok {
		return RunResult{}, ErrRateLimited
	}

	cctx, cancel := context.WithTimeout(ctx, invokeTimeout)
	defer cancel()
	return s.runner.Run(cctx, source, stdin)
}
