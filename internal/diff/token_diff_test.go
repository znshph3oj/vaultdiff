package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseToken = &vault.TokenInfo{
	Accessor:  "acc1",
	Policies:  []string{"default", "read-only"},
	TTL:       3600,
	Renewable: true,
}

func TestCompareTokens_NoChanges(t *testing.T) {
	b := &vault.TokenInfo{
		Accessor:  "acc1",
		Policies:  []string{"default", "read-only"},
		TTL:       3600,
		Renewable: true,
	}
	d := diff.CompareTokens(baseToken, b)
	if len(d.PoliciesAdded) != 0 || len(d.PoliciesRemoved) != 0 {
		t.Error("expected no policy changes")
	}
	if d.TTLChanged {
		t.Error("expected no TTL change")
	}
	if d.RenewableChanged {
		t.Error("expected no renewable change")
	}
}

func TestCompareTokens_PolicyAdded(t *testing.T) {
	b := &vault.TokenInfo{
		Policies:  []string{"default", "read-only", "admin"},
		TTL:       3600,
		Renewable: true,
	}
	d := diff.CompareTokens(baseToken, b)
	if len(d.PoliciesAdded) != 1 || d.PoliciesAdded[0] != "admin" {
		t.Errorf("expected admin added, got %v", d.PoliciesAdded)
	}
}

func TestCompareTokens_PolicyRemoved(t *testing.T) {
	b := &vault.TokenInfo{
		Policies:  []string{"default"},
		TTL:       3600,
		Renewable: true,
	}
	d := diff.CompareTokens(baseToken, b)
	if len(d.PoliciesRemoved) != 1 || d.PoliciesRemoved[0] != "read-only" {
		t.Errorf("expected read-only removed, got %v", d.PoliciesRemoved)
	}
}

func TestCompareTokens_TTLChanged(t *testing.T) {
	b := &vault.TokenInfo{Policies: []string{"default", "read-only"}, TTL: 7200, Renewable: true}
	d := diff.CompareTokens(baseToken, b)
	if !d.TTLChanged {
		t.Fatal("expected TTL change")
	}
	if d.OldTTL != 3600 || d.NewTTL != 7200 {
		t.Errorf("wrong TTL values: %d -> %d", d.OldTTL, d.NewTTL)
	}
}

func TestCompareTokens_RenewableChanged(t *testing.T) {
	b := &vault.TokenInfo{Policies: []string{"default", "read-only"}, TTL: 3600, Renewable: false}
	d := diff.CompareTokens(baseToken, b)
	if !d.RenewableChanged {
		t.Fatal("expected renewable change")
	}
	if d.OldRenewable != true || d.NewRenewable != false {
		t.Errorf("wrong renewable values")
	}
}

func TestFprintTokenDiff_NoChanges(t *testing.T) {
	d := &diff.TokenDiff{}
	var buf bytes.Buffer
	diff.FprintTokenDiff(&buf, d)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}

func TestFprintTokenDiff_ShowsChanges(t *testing.T) {
	d := &diff.TokenDiff{
		PoliciesAdded:   []string{"admin"},
		PoliciesRemoved: []string{"read-only"},
		TTLChanged:      true,
		OldTTL:          3600,
		NewTTL:          7200,
	}
	var buf bytes.Buffer
	diff.FprintTokenDiff(&buf, d)
	out := buf.String()
	if !strings.Contains(out, "admin") {
		t.Error("expected admin in output")
	}
	if !strings.Contains(out, "read-only") {
		t.Error("expected read-only in output")
	}
	if !strings.Contains(out, "3600") {
		t.Error("expected old TTL in output")
	}
}
