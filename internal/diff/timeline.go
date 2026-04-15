package diff

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// TimelineEntry records the diff between two consecutive secret versions.
type TimelineEntry struct {
	FromVersion int
	ToVersion   int
	At          time.Time
	Changes     []Change
}

// Timeline is an ordered slice of TimelineEntry values.
type Timeline []TimelineEntry

// SecretReader can retrieve a secret version's data.
type SecretReader interface {
	GetSecretVersion(ctx context.Context, path string, version int) (map[string]string, error)
	GetSecretMetadata(ctx context.Context, path string) (*vault.SecretMetadata, error)
}

// BuildTimeline constructs a full diff timeline across all available versions
// of a secret path, from oldest to current.
func BuildTimeline(ctx context.Context, r SecretReader, path string, maskKeys []string) (Timeline, error) {
	meta, err := r.GetSecretMetadata(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("fetching metadata for %q: %w", path, err)
	}

	versionNums := make([]int, 0, len(meta.Versions))
	for v := range meta.Versions {
		versionNums = append(versionNums, v)
	}
	sort.Ints(versionNums)

	if len(versionNums) < 2 {
		return Timeline{}, nil
	}

	var timeline Timeline
	for i := 1; i < len(versionNums); i++ {
		from := versionNums[i-1]
		to := versionNums[i]

		oldData, err := r.GetSecretVersion(ctx, path, from)
		if err != nil {
			continue
		}
		newData, err := r.GetSecretVersion(ctx, path, to)
		if err != nil {
			continue
		}

		changes := Compare(oldData, newData, maskKeys)
		entry := TimelineEntry{
			FromVersion: from,
			ToVersion:   to,
			Changes:     changes,
		}
		if vm, ok := meta.Versions[to]; ok {
			entry.At = vm.CreatedTime
		}
		timeline = append(timeline, entry)
	}

	return timeline, nil
}

// PrintTimeline writes a human-readable timeline to w.
func PrintTimeline(w io.Writer, t Timeline, path string) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Timeline for secret: %s\n", path)
	fmt.Fprintf(w, "%s\n", repeatChar('─', 60))
	for _, entry := range t {
		ts := entry.At.Format(time.RFC3339)
		if entry.At.IsZero() {
			ts = "unknown"
		}
		fmt.Fprintf(w, "  v%d → v%d  [%s]  %d change(s)\n",
			entry.FromVersion, entry.ToVersion, ts, len(entry.Changes))
		for _, c := range entry.Changes {
			fmt.Fprintf(w, "    %s\n", Render(c))
		}
	}
}

func repeatChar(ch rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
