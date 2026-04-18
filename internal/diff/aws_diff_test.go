package diff

import (
	"bytes"
	"testing"

	"github.com/user/vaultdiff/internal/vault"
)

var baseAWSRole = &vault.AWSRoleInfo{
	CredentialType: "iam_user",
	RoleARNs:       "arn:aws:iam::123:role/foo",
	PolicyARNs:     "arn:aws:iam::aws:policy/ReadOnly",
	DefaultSTSTTL:  3600,
	MaxSTSTTL:      7200,
}

func TestCompareAWSRoles_NoChanges(t *testing.T) {
	a := *baseAWSRole
	b := *baseAWSRole
	changes := CompareAWSRoles(&a, &b)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestCompareAWSRoles_CredentialTypeChanged(t *testing.T) {
	a := *baseAWSRole
	b := *baseAWSRole
	b.CredentialType = "assumed_role"
	changes := CompareAWSRoles(&a, &b)
	if len(changes) != 1 || changes[0].Field != "credential_type" {
		t.Fatalf("expected credential_type change, got %+v", changes)
	}
}

func TestCompareAWSRoles_TTLChanged(t *testing.T) {
	a := *baseAWSRole
	b := *baseAWSRole
	b.DefaultSTSTTL = 1800
	changes := CompareAWSRoles(&a, &b)
	if len(changes) != 1 || changes[0].Field != "default_sts_ttl" {
		t.Fatalf("expected default_sts_ttl change, got %+v", changes)
	}
}

func TestCompareAWSRoles_NilInputs(t *testing.T) {
	changes := CompareAWSRoles(nil, nil)
	if changes != nil {
		t.Fatal("expected nil for nil inputs")
	}
}

func TestFprintAWSDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintAWSDiff(&buf, nil)
	if buf.String() == "" {
		t.Fatal("expected output for no changes")
	}
}

func TestFprintAWSDiff_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	changes := []AWSDiffEntry{{Field: "credential_type", OldValue: "iam_user", NewValue: "assumed_role"}}
	FprintAWSDiff(&buf, changes)
	out := buf.String()
	if out == "" {
		t.Fatal("expected diff output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("credential_type")) {
		t.Error("expected credential_type in output")
	}
}
