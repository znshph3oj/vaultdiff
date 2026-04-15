package diff

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockFetcher struct {
	mu       sync.Mutex
	versions map[int]map[string]string
}

func (m *mockFetcher) GetSecretVersion(path string, version int) (map[string]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.versions[version]
	if !ok {
		return nil, errors.New("version not found")
	}
	return v, nil
}

func TestWatcher_DetectsChange(t *testing.T) {
	fetcher := &mockFetcher{
		versions: map[int]map[string]string{
			1: {"key": "old"},
			2: {"key": "new"},
		},
	}

	var mu sync.Mutex
	var reports []*Report

	opts := WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(r *Report) {
			mu.Lock()
			reports = append(reports, r)
			mu.Unlock()
		},
	}

	w := NewWatcher(fetcher, opts)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = w.Watch(ctx, "secret/data/test", 1)

	mu.Lock()
	defer mu.Unlock()
	if len(reports) == 0 {
		t.Fatal("expected at least one change report")
	}
	if reports[0].FromVersion != 1 || reports[0].ToVersion != 2 {
		t.Errorf("unexpected versions: from=%d to=%d", reports[0].FromVersion, reports[0].ToVersion)
	}
}

func TestWatcher_NoChangeWhenVersionUnavailable(t *testing.T) {
	fetcher := &mockFetcher{
		versions: map[int]map[string]string{
			1: {"key": "value"},
		},
	}

	called := false
	opts := WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(r *Report) { called = true },
	}

	w := NewWatcher(fetcher, opts)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = w.Watch(ctx, "secret/data/test", 1)
	if called {
		t.Error("OnChange should not be called when next version is unavailable")
	}
}

func TestNewWatcher_DefaultInterval(t *testing.T) {
	w := NewWatcher(&mockFetcher{}, WatchOptions{})
	if w.opts.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", w.opts.Interval)
	}
}
