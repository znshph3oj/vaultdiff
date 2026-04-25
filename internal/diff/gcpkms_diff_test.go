package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baseGCPKMS = &vault.GCPKMSKeyInfo{
	Name:            "my-key",
	KeyRing:         "my-ring",
	CryptoKey:       "my-crypto-key",
	Algorithm:       "GOOGLE_SYMMETRIC_ENCRYPTION",
	ProtectionLevel: "SOFTWARE",
	RotationPeriod:  "7776000s",
	MinVersion:      1,
}

func TestCompareGCPKMSKeys_NoChanges(t *testing.T) {
	copy := *baseGCPKMS
	changes := CompareGCPKMSKeys(baseGCPKMS, &copy)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompareGCPKMSKeys_AlgorithmChanged(t *testing.T) {
	b := *baseGCPKMS
	b.Algorithm = "RSA_SIGN_PKCS1_2048_SHA256"
	changes := CompareGCPKMSKeys(baseGCPKMS, &b)
	if len(changes) != 1 || changes[0].Field != "algorithm" {
		t.Errorf("expected algorithm change, got %+v", changes)
	}
}

func TestCompareGCPKMSKeys_ProtectionLevelChanged(t *testing.T) {
	b := *baseGCPKMS
	b.ProtectionLevel = "HSM"
	changes := CompareGCPKMSKeys(baseGCPKMS, &b)
	if len(changes) != 1 || changes[0].Field != "protection_level" {
		t.Errorf("expected protection_level change, got %+v", changes)
	}
}

func TestCompareGCPKMSKeys_MinVersionChanged(t *testing.T) {
	b := *baseGCPKMS
	b.MinVersion = 3
	changes := CompareGCPKMSKeys(baseGCPKMS, &b)
	if len(changes) != 1 || changes[0].Field != "min_version" {
		t.Errorf("expected min_version change, got %+v", changes)
	}
}

func TestCompareGCPKMSKeys_NilInputs(t *testing.T) {
	if changes := CompareGCPKMSKeys(nil, baseGCPKMS); changes != nil {
		t.Error("expected nil for nil first arg")
	}
	if changes := CompareGCPKMSKeys(baseGCPKMS, nil); changes != nil {
		t.Error("expected nil for nil second arg")
	}
}

func TestFprintGCPKMSDiff_NoChanges(t *testing.T) {
	copy := *baseGCPKMS
	var buf bytes.Buffer
	FprintGCPKMSDiff(&buf, baseGCPKMS, &copy)
	if !strings.Contains(buf.String(), "No GCP KMS key changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestFprintGCPKMSDiff_ShowsChanges(t *testing.T) {
	b := *baseGCPKMS
	b.RotationPeriod = "2592000s"
	var buf bytes.Buffer
	FprintGCPKMSDiff(&buf, baseGCPKMS, &b)
	if !strings.Contains(buf.String(), "rotation_period") {
		t.Errorf("expected rotation_period in output, got: %s", buf.String())
	}
}
