package diff

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// RollbackRequest describes a request to roll back a secret to a prior version.
type RollbackRequest struct {
	Path       string
	TargetVersion int
	DryRun     bool
	RequestedBy string
}

// RollbackResult captures the outcome of a rollback operation.
type RollbackResult struct {
	Path          string
	FromVersion   int
	ToVersion     int
	DryRun        bool
	PerformedAt   time.Time
	Changes       []Change
}

// SecretVersionReader fetches a secret at a specific version.
type SecretVersionReader interface {
	GetSecretVersion(ctx context.Context, path string, version int) (map[string]interface{}, error)
}

// SecretVersionWriter writes (restores) a secret version.
type SecretVersionWriter interface {
	WriteSecret(ctx context.Context, path string, data map[string]interface{}) error
}

// RollbackClient combines read and write capabilities.
type RollbackClient interface {
	SecretVersionReader
	SecretVersionWriter
}

// Rollback performs (or simulates) a rollback of a secret to the target version.
func Rollback(ctx context.Context, client RollbackClient, req RollbackRequest, currentVersion int) (*RollbackResult, error) {
	current, err := client.GetSecretVersion(ctx, req.Path, currentVersion)
	if err != nil {
		return nil, fmt.Errorf("rollback: failed to read current version %d: %w", currentVersion, err)
	}

	target, err := client.GetSecretVersion(ctx, req.Path, req.TargetVersion)
	if err != nil {
		return nil, fmt.Errorf("rollback: failed to read target version %d: %w", req.TargetVersion, err)
	}

	changes := Compare(current, target, false)

	result := &RollbackResult{
		Path:        req.Path,
		FromVersion: currentVersion,
		ToVersion:   req.TargetVersion,
		DryRun:      req.DryRun,
		PerformedAt: time.Now().UTC(),
		Changes:     changes,
	}

	if !req.DryRun {
		if err := client.WriteSecret(ctx, req.Path, target); err != nil {
			return nil, fmt.Errorf("rollback: failed to write secret: %w", err)
		}
	}

	return result, nil
}

// PrintRollbackResult writes a human-readable summary of the rollback to stdout.
func PrintRollbackResult(r *RollbackResult) {
	FprintRollbackResult(os.Stdout, r)
}

// FprintRollbackResult writes a human-readable summary of the rollback to w.
func FprintRollbackResult(w io.Writer, r *RollbackResult) {
	mode := "APPLIED"
	if r.DryRun {
		mode = "DRY RUN"
	}
	fmt.Fprintf(w, "Rollback [%s] %s: v%d → v%d at %s\n",
		mode, r.Path, r.FromVersion, r.ToVersion, r.PerformedAt.Format(time.RFC3339))
	if len(r.Changes) == 0 {
		fmt.Fprintln(w, "  No changes detected.")
		return
	}
	Render(w, r.Changes)
}
