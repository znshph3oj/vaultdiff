package diff

import (
	"bytes"
	"testing"

	"github.com/user/vaultdiff/internal/vault"
)

var basePKICert = &vault.PKICertInfo{
	SerialNumber: "12:34:56",
	CommonName:   "example.com",
	Issuer:       "MyCA",
	NotBefore:    "2024-01-01",
	NotAfter:     "2025-01-01",
	Revoked:      false,
	Mount:        "pki",
}

func TestComparePKICerts_NoChanges(t *testing.T) {
	b := *basePKICert
	changes := ComparePKICerts(basePKICert, &b)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestComparePKICerts_CommonNameChanged(t *testing.T) {
	b := *basePKICert
	b.CommonName = "other.com"
	changes := ComparePKICerts(basePKICert, &b)
	if len(changes) != 1 || changes[0].Field != "CommonName" {
		t.Errorf("expected CommonName change, got %+v", changes)
	}
}

func TestComparePKICerts_Revoked(t *testing.T) {
	b := *basePKICert
	b.Revoked = true
	changes := ComparePKICerts(basePKICert, &b)
	if len(changes) != 1 || changes[0].Field != "Revoked" {
		t.Errorf("expected Revoked change, got %+v", changes)
	}
}

func TestComparePKICerts_NilInputs(t *testing.T) {
	changes := ComparePKICerts(nil, nil)
	if changes != nil {
		t.Errorf("expected nil for nil inputs")
	}
}

func TestFprintPKIDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	FprintPKIDiff(&buf, "12:34:56", nil)
	if !bytes.Contains(buf.Bytes(), []byte("no changes")) {
		t.Errorf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestFprintPKIDiff_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	changes := []PKICertChange{{Field: "Issuer", OldValue: "OldCA", NewValue: "NewCA"}}
	FprintPKIDiff(&buf, "12:34:56", changes)
	if !bytes.Contains(buf.Bytes(), []byte("Issuer")) {
		t.Errorf("expected Issuer in output, got: %s", buf.String())
	}
}
