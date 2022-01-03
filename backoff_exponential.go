package retry

import (
	"context"
	"sync/atomic"
	"time"
)

type exponentialBackoff struct {
	base    time.Duration
	attempt uint64
}

// Exponential is a wrapper around Retry that uses an exponential backoff. It's
// very efficient, but does not check for overflow, so ensure you bound the
// retry. It panics if the given base is less than zero.
func Exponential(ctx context.Context, base time.Duration, f RetryFunc) error {
	return Do(ctx, NewExponential(base), f)
}

// NewExponential creates a new exponential backoff using the starting value of
// base and doubling on each failure (1, 2, 4, 8, 16, 32, 64...), up to max.
// It's very efficient, but does not check for overflow, so ensure you bound the
// retry. It panics if the given base is less than 0.
func NewExponential(base time.Duration) Backoff {
	if base <= 0 {
		panic("base must be greater than 0")
	}

	return &exponentialBackoff{
		base: base,
	}
}

// Next implements Backoff. It is safe for concurrent use.
func (b *exponentialBackoff) Next() (time.Duration, bool) {
	return b.base << (atomic.AddUint64(&b.attempt, 1) - 1), false
}
