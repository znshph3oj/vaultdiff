package diff

import (
	"testing"
)

func TestCompare_AddedKey(t *testing.T) {
	from := map[string]interface{}{"username": "admin"}
	to := map[string]interface{}{"username": "admin", "password": "secret"}

	result := Compare("secret/myapp", 1, 2, from, to)

	if !result.HasChanges() {
		t.Fatal("expected changes but got none")
	}
	if result.FromVersion != 1 || result.ToVersion != 2 {
		t.Errorf("unexpected versions: %d -> %d", result.FromVersion, result.ToVersion)
	}

	var found bool
	for _, c := range result.Changes {
		if c.Key == "password" && c.Type == Added {
			found = true
		}
	}
	if !found {
		t.Error("expected 'password' key to be marked as added")
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	from := map[string]interface{}{"username": "admin", "token": "abc123"}
	to := map[string]interface{}{"username": "admin"}

	result := Compare("secret/myapp", 2, 3, from, to)

	var found bool
	for _, c := range result.Changes {
		if c.Key == "token" && c.Type == Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected 'token' key to be marked as removed")
	}
}

func TestCompare_ModifiedKey(t *testing.T) {
	from := map[string]interface{}{"password": "old"}
	to := map[string]interface{}{"password": "new"}

	result := Compare("secret/myapp", 1, 2, from, to)

	for _, c := range result.Changes {
		if c.Key == "password" {
			if c.Type != Modified {
				t.Errorf("expected Modified, got %s", c.Type)
			}
			if c.OldValue != "old" || c.NewValue != "new" {
				t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
			}
			return
		}
	}
	t.Error("'password' key not found in changes")
}

func TestCompare_UnchangedKey(t *testing.T) {
	from := map[string]interface{}{"host": "localhost"}
	to := map[string]interface{}{"host": "localhost"}

	result := Compare("secret/db", 1, 2, from, to)

	if result.HasChanges() {
		t.Error("expected no changes")
	}
	for _, c := range result.Changes {
		if c.Key == "host" && c.Type != Unchanged {
			t.Errorf("expected Unchanged, got %s", c.Type)
		}
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	from := map[string]interface{}{}
	to := map[string]interface{}{}

	result := Compare("secret/empty", 1, 2, from, to)

	if result.HasChanges() {
		t.Error("expected no changes for two empty maps")
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
}

func TestMaskValue(t *testing.T) {
	masked := MaskValue("supersecret")
	if masked != "********" {
		t.Errorf("expected 8 asterisks, got %q", masked)
	}
	if MaskValue("") != "" {
		t.Error("expected empty string for empty input")
	}
}
