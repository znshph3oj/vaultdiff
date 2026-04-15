package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/vaultdiff/internal/vault"
)

func baseLeaseInfo(id string) *vault.LeaseInfo {
	return &vault.LeaseInfo{
		LeaseID:       id,
		Renewable:     true,
		LeaseDuration: 300 * time.Second,
		ExpireTime:    time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestCompareLeases_NoChanges(t *testing.T) {
	old := baseLeaseInfo("db/creds/role/abc")
	cur := baseLeaseInfo("db/creds/role/abc")

	changes := CompareLeases(old, cur)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestCompareLeases_DurationChanged(t *testing.T) {
	old := baseLeaseInfo("db/creds/role/abc")
	cur := baseLeaseInfo("db/creds/role/abc")
	cur.LeaseDuration = 600 * time.Second

	changes := CompareLeases(old, cur)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "lease_duration" {
		t.Errorf("expected field lease_duration, got %q", changes[0].Field)
	}
}

func TestCompareLeases_RenewableChanged(t *testing.T) {
	old := baseLeaseInfo("db/creds/role/abc")
	cur := baseLeaseInfo("db/creds/role/abc")
	cur.Renewable = false

	changes := CompareLeases(old, cur)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "renewable" {
		t.Errorf("expected field renewable, got %q", changes[0].Field)
	}
}

func TestCompareLeases_NilInputs(t *testing.T) {
	if changes := CompareLeases(nil, nil); changes != nil {
		t.Error("expected nil for nil inputs")
	}
	if changes := CompareLeases(baseLeaseInfo("x"), nil); changes != nil {
		t.Error("expected nil when current is nil")
	}
}

func TestFprintLeaseDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintLeaseDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No lease changes") {
		t.Errorf("expected no-change message, got %q", buf.String())
	}
}

func TestFprintLeaseDiff_WithChanges(t *testing.T) {
	changes := []LeaseChange{
		{LeaseID: "db/creds/role/abc", Field: "lease_duration", OldValue: "5m0s", NewValue: "10m0s"},
	}
	var buf bytes.Buffer
	FprintLeaseDiff(&buf, changes)
	out := buf.String()
	if !strings.Contains(out, "lease_duration") {
		t.Errorf("expected field name in output, got %q", out)
	}
	if !strings.Contains(out, "db/creds/role/abc") {
		t.Errorf("expected lease ID in output, got %q", out)
	}
}
