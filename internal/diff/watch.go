package diff

import (
	"context"
	"fmt"
	"time"
)

// WatchOptions configures the watch behavior.
type WatchOptions struct {
	Interval  time.Duration
	MaskKeys  []string
	OnChange  func(report *Report)
}

// Watcher polls two secret versions and emits a Report when changes are detected.
type Watcher struct {
	opts    WatchOptions
	fetcher SecretFetcher
}

// SecretFetcher abstracts fetching a secret version by path and version number.
type SecretFetcher interface {
	GetSecretVersion(path string, version int) (map[string]string, error)
}

// NewWatcher creates a Watcher with the given fetcher and options.
func NewWatcher(fetcher SecretFetcher, opts WatchOptions) *Watcher {
	if opts.Interval == 0 {
		opts.Interval = 30 * time.Second
	}
	return &Watcher{opts: opts, fetcher: fetcher}
}

// Watch polls the secret at path starting from baseVersion and calls OnChange
// whenever the next version differs from the previous one.
func (w *Watcher) Watch(ctx context.Context, path string, baseVersion int) error {
	prev, err := w.fetcher.GetSecretVersion(path, baseVersion)
	if err != nil {
		return fmt.Errorf("watch: fetch base version %d: %w", baseVersion, err)
	}

	current := baseVersion
	ticker := time.NewTicker(w.opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			next, err := w.fetcher.GetSecretVersion(path, current+1)
			if err != nil {
				// version not yet available; keep polling
				continue
			}
			changes := Compare(prev, next, w.opts.MaskKeys)
			report := NewReport(path, current, current+1, changes)
			if w.opts.OnChange != nil {
				w.opts.OnChange(report)
			}
			prev = next
			current++
		}
	}
}
