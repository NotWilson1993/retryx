# retryx

Small, dependency-free retry helper for Go.

## Install

```bash
go get github.com/NotWilson1993/retryx
```

## Usage

```go
cfg := retryx.Default()
res := retryx.Do(ctx, cfg, func(ctx context.Context) error {
	// do work
	return nil
})

if res.LastErr != nil {
	// all attempts failed
}
```

## Examples

### Retry a database call

```go
cfg := retryx.Default()
res := retryx.Do(ctx, cfg, func(ctx context.Context) error {
	return db.PingContext(ctx)
})

if res.LastErr != nil {
	return res.LastErr
}
```

### Retry a queue publish

```go
cfg := retryx.Default()
cfg.Attempts = 5
cfg.Base = 100 * time.Millisecond

res := retryx.Do(ctx, cfg, func(ctx context.Context) error {
	return queue.Publish(ctx, msg)
})

if res.LastErr != nil {
	log.Printf("publish failed after %d attempts: %v", res.Attempts, res.LastErr)
}
```

### Retry with jitter

```go
cfg := retryx.Default()
cfg.Jitter = 250 * time.Millisecond

_ = retryx.Do(ctx, cfg, func(ctx context.Context) error {
	return doWork(ctx)
})
```

## Notes

- Retries stop on context cancellation.
- Exponential backoff with optional jitter.
- No global state and no magic.
