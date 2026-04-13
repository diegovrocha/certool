package inspect

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractCN(t *testing.T) {
	tests := []struct {
		subject  string
		expected string
	}{
		{"C=BR, O=ICP-Brasil, CN=Meu Certificado", "Meu Certificado"},
		{"CN=meusite.com.br", "meusite.com.br"},
		{"DC=ihf, DC=corp, CN=Itau-R0650", "Itau-R0650"},
		{"C=BR, O=Test", "C=BR, O=Test"}, // no CN returns full subject
		{"CN=RV2 SERVICOS, OU=Certificado PJ A1", "RV2 SERVICOS"},
	}

	for _, tt := range tests {
		result := extractCN(tt.subject)
		if result != tt.expected {
			t.Errorf("extractCN(%q) = %q, expected %q", tt.subject, result, tt.expected)
		}
	}
}

func TestSplitPEM(t *testing.T) {
	// Generate a test cert
	dir := t.TempDir()
	certFile := filepath.Join(dir, "test.crt")
	keyFile := filepath.Join(dir, "test.key")

	err := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-nodes",
		"-keyout", keyFile, "-out", certFile,
		"-days", "1", "-subj", "/CN=test").Run()
	if err != nil {
		t.Fatalf("Failed to generate test cert: %v", err)
	}

	// 1 cert
	certs := splitPEM(certFile)
	if len(certs) != 1 {
		t.Errorf("splitPEM with 1 cert should return 1, returned %d", len(certs))
	}

	// 3 certs (simulated chain)
	chainFile := filepath.Join(dir, "chain.pem")
	certData, _ := os.ReadFile(certFile)
	chainData := string(certData) + "\n" + string(certData) + "\n" + string(certData)
	os.WriteFile(chainFile, []byte(chainData), 0644)

	certs = splitPEM(chainFile)
	if len(certs) != 3 {
		t.Errorf("splitPEM with 3 certs should return 3, returned %d", len(certs))
	}

	// Clean up tmp files
	for _, c := range certs {
		os.Remove(c)
	}
}

func TestExtractInfo(t *testing.T) {
	dir := t.TempDir()
	certFile := filepath.Join(dir, "test.crt")
	keyFile := filepath.Join(dir, "test.key")

	err := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-nodes",
		"-keyout", keyFile, "-out", certFile,
		"-days", "365", "-subj", "/CN=teste.local/O=TestOrg/C=BR").Run()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	info := extractInfo(certFile, 1, 1)
	if info == nil {
		t.Fatal("extractInfo returned nil")
	}

	if info.CN != "teste.local" {
		t.Errorf("CN expected 'teste.local', got '%s'", info.CN)
	}

	if !strings.Contains(info.Subject, "teste.local") {
		t.Errorf("Subject should contain 'teste.local': %s", info.Subject)
	}

	// Self-signed = subject == issuer
	if info.Subject != info.Issuer {
		t.Error("Self-signed cert should have subject == issuer")
	}

	if !strings.Contains(info.CertType, "Root CA") {
		t.Errorf("Self-signed cert should be Root CA, got: %s", info.CertType)
	}

	if info.Valid != "Valid" {
		t.Errorf("Cert with 365 days should be valid, got: %s", info.Valid)
	}

	if info.Serial == "" {
		t.Error("Serial should not be empty")
	}

	if info.Fingerprint == "" {
		t.Error("Fingerprint should not be empty")
	}
}

func TestExtractInfoShortLived(t *testing.T) {
	dir := t.TempDir()
	certFile := filepath.Join(dir, "short.crt")
	keyFile := filepath.Join(dir, "short.key")

	// Cert with 1 day validity (expires soon)
	err := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-nodes",
		"-keyout", keyFile, "-out", certFile,
		"-days", "1", "-subj", "/CN=shortlived").Run()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	info := extractInfo(certFile, 1, 1)
	if info == nil {
		t.Fatal("extractInfo returned nil")
	}

	// 1-day cert should be marked as "expires in less than 30 days"
	if info.Valid == "" {
		t.Error("Valid should not be empty")
	}
	if info.ValidColor != "yellow" {
		t.Errorf("1-day cert should have yellow color, got: %s", info.ValidColor)
	}
}

func TestExtractInfoInvalidFile(t *testing.T) {
	dir := t.TempDir()
	badFile := filepath.Join(dir, "bad.pem")
	os.WriteFile(badFile, []byte("not a certificate"), 0644)

	info := extractInfo(badFile, 1, 1)
	if info != nil {
		t.Error("extractInfo should return nil for invalid file")
	}
}

func TestHasPEMMarker(t *testing.T) {
	dir := t.TempDir()

	// PEM file
	pemFile := filepath.Join(dir, "test.pem")
	os.WriteFile(pemFile, []byte("-----BEGIN CERTIFICATE-----\ndata\n-----END CERTIFICATE-----\n"), 0644)
	if !hasPEMMarker(pemFile) {
		t.Error("hasPEMMarker should return true for PEM")
	}

	// Binary file
	derFile := filepath.Join(dir, "test.der")
	os.WriteFile(derFile, []byte{0x30, 0x82, 0x01}, 0644)
	if hasPEMMarker(derFile) {
		t.Error("hasPEMMarker should return false for DER")
	}
}

func TestHasPrivateKeyOnly(t *testing.T) {
	dir := t.TempDir()

	// File with key + cert
	bothFile := filepath.Join(dir, "both.pem")
	os.WriteFile(bothFile, []byte("-----BEGIN PRIVATE KEY-----\nkey\n-----END PRIVATE KEY-----\n-----BEGIN CERTIFICATE-----\ncert\n-----END CERTIFICATE-----\n"), 0644)
	if hasPrivateKeyOnly(bothFile) {
		t.Error("File with cert+key should not be 'private key only'")
	}

	// File with key only
	keyFile := filepath.Join(dir, "key.pem")
	os.WriteFile(keyFile, []byte("-----BEGIN PRIVATE KEY-----\nkey\n-----END PRIVATE KEY-----\n"), 0644)
	if !hasPrivateKeyOnly(keyFile) {
		t.Error("File with only key should be 'private key only'")
	}
}
