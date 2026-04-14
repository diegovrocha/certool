package batch

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func ensureOpenSSL(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("openssl not installed")
	}
}

func genCert(t *testing.T, dir, name string, days int) string {
	t.Helper()
	certFile := filepath.Join(dir, name+".crt")
	keyFile := filepath.Join(dir, name+".key")
	err := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-nodes",
		"-keyout", keyFile, "-out", certFile,
		"-days", itoa(days), "-subj", "/CN="+name).Run()
	if err != nil {
		t.Fatalf("openssl failed: %v", err)
	}
	return certFile
}

func itoa(n int) string {
	// simple because we expect small int + possible negative
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}

func TestScanDir(t *testing.T) {
	ensureOpenSSL(t)
	dir := t.TempDir()
	genCert(t, dir, "valid", 365)
	genCert(t, dir, "soon", 10)

	rows := scanDir(dir, 5)
	if len(rows) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(rows))
	}

	foundValid := false
	foundSoon := false
	for _, r := range rows {
		if r.CN == "valid" {
			foundValid = true
			if r.Status != "ok" {
				t.Errorf("valid cert should have status 'ok', got %q", r.Status)
			}
			if r.Days < 300 {
				t.Errorf("valid cert should have >300 days, got %d", r.Days)
			}
		}
		if r.CN == "soon" {
			foundSoon = true
			if r.Status != "warn" {
				t.Errorf("soon cert should have status 'warn' (<=30 days), got %q (days=%d)", r.Status, r.Days)
			}
		}
	}
	if !foundValid {
		t.Error("did not find 'valid' cert in scan results")
	}
	if !foundSoon {
		t.Error("did not find 'soon' cert in scan results")
	}
}

func TestSortByDays(t *testing.T) {
	m := &Model{
		sortMode: sortByDays,
		rows: []Row{
			{CN: "c", Days: 30},
			{CN: "a", Days: 100},
			{CN: "b", Days: -5},
		},
	}
	m.applySort()
	if m.rows[0].Days != -5 || m.rows[1].Days != 30 || m.rows[2].Days != 100 {
		t.Errorf("sortByDays: got %v", []int{m.rows[0].Days, m.rows[1].Days, m.rows[2].Days})
	}
}

func TestSortByName(t *testing.T) {
	m := &Model{
		sortMode: sortByCN,
		rows: []Row{
			{CN: "charlie"},
			{CN: "alpha"},
			{CN: "bravo"},
		},
	}
	m.applySort()
	order := []string{m.rows[0].CN, m.rows[1].CN, m.rows[2].CN}
	want := []string{"alpha", "bravo", "charlie"}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("sortByCN order: got %v want %v", order, want)
			break
		}
	}
}

func TestStatusColor(t *testing.T) {
	// Test status classification via formatRow output (uses style markers)
	// but more reliably: exercise extractRow logic which assigns statuses.
	// Since we have direct status strings, we test the thresholds by
	// constructing Rows and verifying formatRow produces expected substrings.
	cases := []struct {
		days   int
		status string
		marker string
	}{
		{40, "ok", "OK"},
		{15, "warn", "15d"},
		{-5, "expired", "EXPIRED"},
	}
	for _, c := range cases {
		r := Row{CN: "x", Days: c.days, Status: c.status, NotAfter: "2026-01-01"}
		out := formatRow(r, 24, 28, 12, 5)
		if !strings.Contains(out, c.marker) {
			t.Errorf("days=%d status=%s expected marker %q in row: %q", c.days, c.status, c.marker, out)
		}
	}
}
