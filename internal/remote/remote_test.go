package remote

import (
	"strings"
	"testing"
)

func TestParseHostPort(t *testing.T) {
	tests := []struct {
		raw      string
		wantHost string
		wantPort string
	}{
		{"example.com", "example.com", "443"},
		{"example.com:8443", "example.com", "8443"},
		{"host:1234", "host", "1234"},
		{"", "", "443"},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			h, p := parseHostPort(tt.raw)
			if h != tt.wantHost || p != tt.wantPort {
				t.Errorf("parseHostPort(%q) = (%q,%q), want (%q,%q)",
					tt.raw, h, p, tt.wantHost, tt.wantPort)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"example.com", "example.com"},
		{"a/b", "a_b"},
		{"a:b", "a_b"},
		{"*.example.com", "_.example.com"},
	}
	for _, tt := range tests {
		if got := sanitizeName(tt.in); got != tt.want {
			t.Errorf("sanitizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestDefaultSaveName(t *testing.T) {
	if got := defaultSaveName("example.com"); got != "example.com_chain.pem" {
		t.Errorf("defaultSaveName: got %q", got)
	}
	if got := defaultSaveName(""); got != "remote_chain.pem" {
		t.Errorf("defaultSaveName empty: got %q", got)
	}
}

func TestExtractCertBlocks(t *testing.T) {
	raw := `random prefix
-----BEGIN CERTIFICATE-----
AAA
-----END CERTIFICATE-----
middle
-----BEGIN CERTIFICATE-----
BBB
-----END CERTIFICATE-----
end`
	blocks := extractCertBlocks(raw)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
}

func TestFirstErrorLine(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "connect refused",
			raw:  "CONNECTED(00000003)\nconnect:errno=61\nConnection refused\n",
			want: "connect:errno=61",
		},
		{
			name: "dns failure",
			raw:  "getaddrinfo: nodename nor servname provided, or not known\n",
			want: "getaddrinfo: nodename nor servname provided, or not known",
		},
		{
			name: "verify error",
			raw:  "depth=0 CN = example.com\nverify error:num=20:unable to get local issuer certificate\nverify return:1\n",
			want: "verify error:num=20:unable to get local issuer certificate",
		},
		{
			name: "ssl routine",
			raw:  "40AB4E01FF7F0000:error:0A00010B:SSL routines:ssl3_get_record:wrong version number:\n",
			want: "40AB4E01FF7F0000:error:0A00010B:SSL routines:ssl3_get_record:wrong version number:",
		},
		{
			name: "no diagnostic",
			raw:  "just some noise\nwith nothing useful\n",
			want: "",
		},
		{
			name: "empty",
			raw:  "",
			want: "",
		},
		{
			name: "truncates long messages",
			raw:  "SSL routines:" + strings.Repeat("x", 200),
			want: "SSL routines:" + strings.Repeat("x", 117-len("SSL routines:")) + "...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstErrorLine(tt.raw)
			if got != tt.want {
				t.Errorf("firstErrorLine(%q)\n  got:  %q\n  want: %q", tt.name, got, tt.want)
			}
		})
	}
}
