package utils

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

type RetryWithBackoffParam struct {
	Func       func() error
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func RetryWithBackoff(ctx context.Context, param RetryWithBackoffParam) error {
	var lastErr error

	for attempt := 0; attempt <= param.MaxRetries; attempt++ {
		if err := param.Func(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		backoff := float64(param.BaseDelay) * math.Pow(2, float64(attempt))
		if backoff > float64(param.MaxDelay) {
			backoff = float64(param.MaxDelay)
		}

		jitter := time.Duration(backoff * (0.5 + rand.Float64()))
		select {
		case <-time.After(jitter):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if lastErr == nil {
		lastErr = errors.New("unknown error occurred on exp backoff")
	}
	return lastErr
}
