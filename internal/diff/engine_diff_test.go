package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

var kvEngine = vault.EngineInfo{Path: "secret/", Type: "kv", Description: "kv store"}
var pkiEngine = vault.EngineInfo{Path: "pki/", Type: "pki", Description: "pki"}
var transitEngine = vault.EngineInfo{Path: "transit/", Type: "transit", Description: "transit"}

func TestCompareEngines_Added(t *testing.T) {
	entries := diff.CompareEngines([]vault.EngineInfo{kvEngine}, []vault.EngineInfo{kvEngine, pkiEngine})
	if len(entries) != 1 || entries[0].Status != "added" || entries[0].Path != "pki/" {
		t.Fatalf("expected 1 added entry for pki/, got %+v", entries)
	}
}

func TestCompareEngines_Removed(t *testing.T) {
	entries := diff.CompareEngines([]vault.EngineInfo{kvEngine, pkiEngine}, []vault.EngineInfo{kvEngine})
	if len(entries) != 1 || entries[0].Status != "removed" || entries[0].Path != "pki/" {
		t.Fatalf("expected 1 removed entry for pki/, got %+v", entries)
	}
}

func TestCompareEngines_Modified(t *testing.T) {
	modified := vault.EngineInfo{Path: "secret/", Type: "kv-v2", Description: "upgraded"}
	entries := diff.CompareEngines([]vault.EngineInfo{kvEngine}, []vault.EngineInfo{modified})
	if len(entries) != 1 || entries[0].Status != "modified" {
		t.Fatalf("expected 1 modified entry, got %+v", entries)
	}
}

func TestCompareEngines_NoChanges(t *testing.T) {
	entries := diff.CompareEngines([]vault.EngineInfo{kvEngine}, []vault.EngineInfo{kvEngine})
	if len(entries) != 0 {
		t.Fatalf("expected no changes, got %+v", entries)
	}
}

func TestFprintEngineDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	diff.FprintEngineDiff(&buf, nil)
	if !strings.Contains(buf.String(), "No engine changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestFprintEngineDiff_ShowsEntries(t *testing.T) {
	entries := diff.CompareEngines(
		[]vault.EngineInfo{kvEngine},
		[]vault.EngineInfo{kvEngine, transitEngine},
	)
	var buf bytes.Buffer
	diff.FprintEngineDiff(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "transit/") {
		t.Errorf("expected transit/ in output, got: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' marker in output, got: %s", out)
	}
}
