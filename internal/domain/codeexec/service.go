package codeexec

import (
	"context"
	"strings"

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

func (s *Service) GradePreview(ctx context.Context, userID uuid.UUID, source string, tests []PreviewTest) (int, []PreviewResult, error) {
	if len(source) == 0 {
		return 0, nil, ErrInvalid
	}
	if len(source) > maxSourceBytes || len(tests) > maxPreviewTests {
		return 0, nil, ErrTooLarge
	}
	for _, t := range tests {
		if len(t.Stdin) > maxStdinBytes {
			return 0, nil, ErrTooLarge
		}
	}
	if ok, err := s.redis.AcquireLock(ctx, "coderun:"+userID.String(), minRunInterval); err == nil && !ok {
		return 0, nil, ErrRateLimited
	}

	inputs := make([]string, len(tests))
	for i, t := range tests {
		inputs[i] = t.Stdin
	}

	gctx, cancel := context.WithTimeout(ctx, gradePreviewTimeout)
	defer cancel()
	runs, err := s.runner.Grade(gctx, source, inputs)
	if err != nil {
		return 0, nil, err
	}

	passed := 0
	results := make([]PreviewResult, 0, len(tests))
	for i, t := range tests {
		var run RunResult
		if i < len(runs) {
			run = runs[i]
		}
		ok := !run.TimedOut && normalizeOutput(run.Stdout) == normalizeOutput(t.ExpectedStdout)
		if ok {
			passed++
		}
		results = append(results, PreviewResult{
			Ordinal:        i,
			Passed:         ok,
			TimedOut:       run.TimedOut,
			ActualStdout:   run.Stdout,
			ExpectedStdout: t.ExpectedStdout,
			Stderr:         run.Stderr,
			DurationMs:     run.DurationMs,
		})
	}
	return passed, results, nil
}

func normalizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}
