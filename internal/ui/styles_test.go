package ui

import (
	"strings"
	"testing"
)

func TestBanner(t *testing.T) {
	b := Banner()
	// Banner uses lipgloss which may or may not include ANSI codes
	// depending on terminal detection. Check the raw ASCII text.
	if !strings.Contains(b, "___") {
		t.Error("Banner should contain ASCII art")
	}
	if !strings.Contains(b, Version) {
		t.Errorf("Banner should contain version '%s'", Version)
	}
	if !strings.Contains(b, "certificate") {
		t.Error("Banner should contain subtitle")
	}
}

func TestOpenSSLVersion(t *testing.T) {
	v := OpenSSLVersion()
	if v == "" {
		t.Error("OpenSSLVersion should not be empty")
	}
	if strings.Contains(v, "not found") {
		t.Skip("OpenSSL not installed")
	}
	if !strings.Contains(strings.ToLower(v), "openssl") {
		t.Errorf("OpenSSLVersion should contain 'openssl', got: %s", v)
	}
}

func TestResultBox(t *testing.T) {
	ok := ResultBox(true, "Success", "file.pem")
	if !strings.Contains(ok, "Success") {
		t.Error("ResultBox success should contain title")
	}

	fail := ResultBox(false, "Error", "message")
	if !strings.Contains(fail, "Error") {
		t.Error("ResultBox error should contain title")
	}
}
