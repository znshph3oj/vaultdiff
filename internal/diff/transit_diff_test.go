package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseTransitKey = &vault.TransitKeyInfo{
	Name:              "mykey",
	Type:              "aes256-gcm96",
	DeletionAllowed:   false,
	Exportable:        false,
	LatestVersion:     1,
	MinDecryptVersion: 1,
}

func TestCompareTransitKeys_NoChanges(t *testing.T) {
	changes := CompareTransitKeys(baseTransitKey, baseTransitKey)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareTransitKeys_TypeChanged(t *testing.T) {
	b := *baseTransitKey
	b.Type = "rsa-2048"
	changes := CompareTransitKeys(baseTransitKey, &b)
	if len(changes) != 1 || changes[0].Field != "type" {
		t.Errorf("expected type change, got %+v", changes)
	}
}

func TestCompareTransitKeys_VersionChanged(t *testing.T) {
	b := *baseTransitKey
	b.LatestVersion = 5
	changes := CompareTransitKeys(baseTransitKey, &b)
	if len(changes) != 1 || changes[0].Field != "latest_version" {
		t.Errorf("expected latest_version change, got %+v", changes)
	}
}

func TestCompareTransitKeys_NilInputs(t *testing.T) {
	if changes := CompareTransitKeys(nil, baseTransitKey); changes != nil {
		t.Error("expected nil for nil input")
	}
	if changes := CompareTransitKeys(baseTransitKey, nil); changes != nil {
		t.Error("expected nil for nil input")
	}
}

func TestFprintTransitDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintTransitDiff(&buf, "mykey", nil)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got %q", buf.String())
	}
}

func TestFprintTransitDiff_WithChanges(t *testing.T) {
	changes := []TransitKeyChange{{Field: "type", Old: "aes256-gcm96", New: "rsa-2048"}}
	var buf bytes.Buffer
	FprintTransitDiff(&buf, "mykey", changes)
	if !strings.Contains(buf.String(), "type") {
		t.Errorf("expected field 'type' in output, got %q", buf.String())
	}
}
