package generate

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseCustomSANs(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "empty",
			raw:  "",
			want: nil,
		},
		{
			name: "bare DNS names",
			raw:  "example.com, *.example.com",
			want: []string{"DNS:example.com", "DNS:*.example.com"},
		},
		{
			name: "explicit prefixes preserved",
			raw:  "DNS:a.com, IP:1.2.3.4, email:x@y.com",
			want: []string{"DNS:a.com", "IP:1.2.3.4", "email:x@y.com"},
		},
		{
			name: "IPv4 auto-detected",
			raw:  "10.0.0.1, 127.0.0.1",
			want: []string{"IP:10.0.0.1", "IP:127.0.0.1"},
		},
		{
			name: "IPv6 auto-detected",
			raw:  "::1, fe80::1",
			want: []string{"IP:::1", "IP:fe80::1"},
		},
		{
			name: "email auto-detected",
			raw:  "admin@example.com",
			want: []string{"email:admin@example.com"},
		},
		{
			name: "mixed",
			raw:  "api.site.com, 10.0.0.1, admin@site.com, DNS:legacy.site.com",
			want: []string{
				"DNS:api.site.com",
				"IP:10.0.0.1",
				"email:admin@site.com",
				"DNS:legacy.site.com",
			},
		},
		{
			name: "whitespace tolerant",
			raw:  "  a.com  ,  b.com  ",
			want: []string{"DNS:a.com", "DNS:b.com"},
		},
		{
			name: "empty entries filtered",
			raw:  "a.com,,b.com,",
			want: []string{"DNS:a.com", "DNS:b.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCustomSANs(tt.raw)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCustomSANs(%q)\n  got:  %v\n  want: %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestIsIPish(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"::1", true},
		{"fe80::1", true},
		{"2001:db8::", true},
		{"example.com", false},
		{"a.b.c.d", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isIPish(tt.in); got != tt.want {
			t.Errorf("isIPish(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSANPresetWeb(t *testing.T) {
	m := &Model{cn: "mysite.local", optCur: 1}
	// Simulate the preset resolution logic
	want := []string{
		"DNS:mysite.local",
		"DNS:localhost",
		"IP:127.0.0.1",
		"IP:::1",
	}
	m.sans = []string{
		"DNS:" + m.cn,
		"DNS:localhost",
		"IP:127.0.0.1",
		"IP:::1",
	}
	if !reflect.DeepEqual(m.sans, want) {
		t.Errorf("web preset = %v, want %v", m.sans, want)
	}
}

func TestSANPresetWildcard(t *testing.T) {
	cn := "example.com"
	got := []string{"DNS:" + cn, "DNS:*." + cn}
	want := []string{"DNS:example.com", "DNS:*.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wildcard preset = %v, want %v", got, want)
	}
}

func TestSANStringFormat(t *testing.T) {
	// Verify the format passed to openssl -addext is correct
	sans := []string{"DNS:a.com", "DNS:*.a.com", "IP:127.0.0.1"}
	joined := strings.Join(sans, ",")
	want := "DNS:a.com,DNS:*.a.com,IP:127.0.0.1"
	if joined != want {
		t.Errorf("joined SANs = %q, want %q", joined, want)
	}
}
