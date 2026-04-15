package diff

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

// fakeRollbackClient implements RollbackClient for testing.
type fakeRollbackClient struct {
	versions    map[int]map[string]interface{}
	written     map[string]interface{}
	writeErr    error
}

func (f *fakeRollbackClient) GetSecretVersion(_ context.Context, _ string, version int) (map[string]interface{}, error) {
	v, ok := f.versions[version]
	if !ok {
		return nil, errors.New("version not found")
	}
	return v, nil
}

func (f *fakeRollbackClient) WriteSecret(_ context.Context, _ string, data map[string]interface{}) error {
	if f.writeErr != nil {
		return f.writeErr
	}
	f.written = data
	return nil
}

func TestRollback_DryRun_DoesNotWrite(t *testing.T) {
	client := &fakeRollbackClient{
		versions: map[int]map[string]interface{}{
			3: {"key": "new"},
			1: {"key": "old"},
		},
	}
	req := RollbackRequest{Path: "secret/app", TargetVersion: 1, DryRun: true}
	result, err := Rollback(context.Background(), client, req, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun to be true")
	}
	if client.written != nil {
		t.Error("expected no write on dry run")
	}
}

func TestRollback_AppliesWrite(t *testing.T) {
	client := &fakeRollbackClient{
		versions: map[int]map[string]interface{}{
			3: {"key": "new"},
			1: {"key": "old"},
		},
	}
	req := RollbackRequest{Path: "secret/app", TargetVersion: 1, DryRun: false}
	_, err := Rollback(context.Background(), client, req, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.written == nil {
		t.Fatal("expected secret to be written")
	}
	if client.written["key"] != "old" {
		t.Errorf("expected written key=old, got %v", client.written["key"])
	}
}

func TestRollback_TargetVersionNotFound(t *testing.T) {
	client := &fakeRollbackClient{
		versions: map[int]map[string]interface{}{
			3: {"key": "new"},
		},
	}
	req := RollbackRequest{Path: "secret/app", TargetVersion: 99, DryRun: true}
	_, err := Rollback(context.Background(), client, req, 3)
	if err == nil {
		t.Fatal("expected error for missing target version")
	}
}

func TestFprintRollbackResult_DryRunLabel(t *testing.T) {
	r := &RollbackResult{
		Path: "secret/app", FromVersion: 3, ToVersion: 1, DryRun: true,
		Changes: []Change{{Key: "key", Type: Modified, OldValue: "new", NewValue: "old"}},
	}
	var buf bytes.Buffer
	FprintRollbackResult(&buf, r)
	if !strings.Contains(buf.String(), "DRY RUN") {
		t.Errorf("expected DRY RUN in output, got: %s", buf.String())
	}
}

func TestFprintRollbackResult_NoChanges(t *testing.T) {
	r := &RollbackResult{
		Path: "secret/app", FromVersion: 2, ToVersion: 2, DryRun: false,
		Changes: nil,
	}
	var buf bytes.Buffer
	FprintRollbackResult(&buf, r)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' in output, got: %s", buf.String())
	}
}
