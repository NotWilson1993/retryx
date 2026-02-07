package retryx

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoSuccess(t *testing.T) {
	calls := 0
	res := Do(context.TODO(), Default(), func(ctx context.Context) error {
		calls++
		return nil
	})
	if res.LastErr != nil {
		t.Fatalf("unexpected error: %v", res.LastErr)
	}
	if res.Attempts != 1 || calls != 1 {
		t.Fatalf("expected 1 attempt, got attempts=%d calls=%d", res.Attempts, calls)
	}
}

func TestDoRetriesAndFails(t *testing.T) {
	calls := 0
	cfg := Default()
	cfg.Attempts = 3
	cfg.Base = 1 * time.Millisecond
	cfg.Max = 2 * time.Millisecond
	res := Do(context.TODO(), cfg, func(ctx context.Context) error {
		calls++
		return errors.New("fail")
	})
	if res.LastErr == nil {
		t.Fatalf("expected error")
	}
	if res.Attempts != 3 || calls != 3 {
		t.Fatalf("expected 3 attempts, got attempts=%d calls=%d", res.Attempts, calls)
	}
}

func TestDoContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	res := Do(ctx, Default(), func(ctx context.Context) error {
		return nil
	})
	if res.LastErr == nil {
		t.Fatalf("expected context error")
	}
	if res.Attempts != 0 {
		t.Fatalf("expected 0 attempts, got %d", res.Attempts)
	}
}

func TestNormalizeDefaults(t *testing.T) {
	cfg := normalize(Config{})
	if cfg.Attempts != 3 {
		t.Fatalf("unexpected Attempts: %d", cfg.Attempts)
	}
	if cfg.Base != 200*time.Millisecond {
		t.Fatalf("unexpected Base: %v", cfg.Base)
	}
	if cfg.Max != 2*time.Second {
		t.Fatalf("unexpected Max: %v", cfg.Max)
	}
}
