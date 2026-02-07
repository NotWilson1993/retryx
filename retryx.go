package retryx

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"
)

// Config controls retry behavior.
type Config struct {
	Attempts int           // default: 3
	Base     time.Duration // default: 200ms
	Max      time.Duration // default: 2s
	Jitter   time.Duration // optional
}

// Default returns the default retry configuration.
func Default() Config {
	return Config{
		Attempts: 3,
		Base:     200 * time.Millisecond,
		Max:      2 * time.Second,
	}
}

// Result describes a retry execution.
type Result struct {
	Attempts int
	LastErr  error
}

// Do retries fn(ctx) until it returns nil or attempts are exhausted.
// It stops immediately if ctx is canceled.
func Do(ctx context.Context, cfg Config, fn func(context.Context) error) Result {
	cfg = normalize(cfg)

	var lastErr error
	for i := 1; i <= cfg.Attempts; i++ {
		if err := ctx.Err(); err != nil {
			return Result{Attempts: i - 1, LastErr: err}
		}

		err := fn(ctx)
		if err == nil {
			return Result{Attempts: i, LastErr: nil}
		}
		lastErr = err

		if i == cfg.Attempts {
			break
		}

		sleep := backoff(cfg, i)
		if sleep > 0 {
			if !sleepWithContext(ctx, sleep) {
				return Result{Attempts: i, LastErr: ctx.Err()}
			}
		}
	}

	return Result{Attempts: cfg.Attempts, LastErr: lastErr}
}

func normalize(cfg Config) Config {
	if cfg.Attempts <= 0 {
		cfg.Attempts = 3
	}
	if cfg.Base <= 0 {
		cfg.Base = 200 * time.Millisecond
	}
	if cfg.Max <= 0 {
		cfg.Max = 2 * time.Second
	}
	if cfg.Max < cfg.Base {
		cfg.Max = cfg.Base
	}
	return cfg
}

func backoff(cfg Config, attempt int) time.Duration {
	d := cfg.Base
	for i := 1; i < attempt; i++ {
		d *= 2
		if d >= cfg.Max {
			d = cfg.Max
			break
		}
	}

	if cfg.Jitter > 0 {
		d += jitter(cfg.Jitter)
	}
	return d
}

func jitter(max time.Duration) time.Duration {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return time.Duration(n.Int64())
}

func sleepWithContext(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
