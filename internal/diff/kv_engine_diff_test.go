package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func baseKVEngine() *vault.KVEngineInfo {
	return &vault.KVEngineInfo{
		Path:    "secret",
		Version: 2,
		Options: map[string]string{"version": "2"},
	}
}

func TestCompareKVEngines_NoChanges(t *testing.T) {
	a, b := baseKVEngine(), baseKVEngine()
	changes := CompareKVEngines(a, b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareKVEngines_VersionChanged(t *testing.T) {
	a := baseKVEngine()
	b := baseKVEngine()
	b.Version = 1
	b.Options["version"] = "1"
	changes := CompareKVEngines(a, b)
	found := false
	for _, c := range changes {
		if c.Field == "version" {
			found = true
		}
	}
	if !found {
		t.Error("expected version change")
	}
}

func TestCompareKVEngines_OptionAdded(t *testing.T) {
	a := baseKVEngine()
	b := baseKVEngine()
	b.Options["max_versions"] = "10"
	changes := CompareKVEngines(a, b)
	if len(changes) == 0 {
		t.Error("expected option change")
	}
}

func TestCompareKVEngines_NilInputs(t *testing.T) {
	changes := CompareKVEngines(nil, nil)
	if len(changes) != 0 {
		t.Errorf("expected no changes for nil inputs")
	}
}

func TestFprintKVEngineDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintKVEngineDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No KV engine changes") {
		t.Error("expected no-changes message")
	}
}

func TestFprintKVEngineDiff_WithChanges(t *testing.T) {
	changes := []KVEngineChange{{Field: "version", OldValue: "1", NewValue: "2"}}
	var buf bytes.Buffer
	FprintKVEngineDiff(&buf, changes)
	if !strings.Contains(buf.String(), "version") {
		t.Error("expected version in output")
	}
}
