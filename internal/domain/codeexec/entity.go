package codeexec

import (
	"errors"
	"time"
)

var (
	ErrNotConfigured = errors.New("code execution is not configured")
	ErrRateLimited   = errors.New("too many runs, slow down")
	ErrTooLarge      = errors.New("code or input too large")
	ErrInvalid       = errors.New("invalid request")
)

const (
	maxSourceBytes = 64 * 1024
	maxStdinBytes  = 16 * 1024
	runTimeoutMs   = 8000
	// minRunInterval throttles a single user to ~one run per window (anti-spam).
	// Lambda reserved concurrency is the hard cost/blast-radius ceiling.
	minRunInterval = 2 * time.Second
	invokeTimeout  = 15 * time.Second
)

// RunResult is what the Lambda returns for one execution.
type RunResult struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exitCode"`
	TimedOut   bool   `json:"timedOut"`
	DurationMs int    `json:"durationMs"`
}
