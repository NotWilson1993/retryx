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

## Notes

- Retries stop on context cancellation.
- Exponential backoff with optional jitter.
- No global state and no magic.
