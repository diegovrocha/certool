package fetchca

import "testing"

func TestParseCAIssuersURL(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "ca issuers with http uri",
			text: `        Authority Information Access:
                CA Issuers - URI:http://pki.example.com/rootCA.crt
                OCSP - URI:http://ocsp.example.com
`,
			want: "http://pki.example.com/rootCA.crt",
		},
		{
			name: "ca issuers with https uri and trailing spaces",
			text: "    CA Issuers - URI:https://ca.example.com/issuer.p7c   \n",
			want: "https://ca.example.com/issuer.p7c",
		},
		{
			name: "only ocsp present",
			text: "        OCSP - URI:http://ocsp.example.com\n",
			want: "",
		},
		{
			name: "no aia at all",
			text: "X509v3 Basic Constraints: CA:FALSE\n",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCAIssuersURL(tt.text); got != tt.want {
				t.Errorf("parseCAIssuersURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCertBlocks(t *testing.T) {
	raw := `subject=CN=Intermediate CA
-----BEGIN CERTIFICATE-----
AAAA
-----END CERTIFICATE-----
subject=CN=Root CA
-----BEGIN CERTIFICATE-----
BBBB
-----END CERTIFICATE-----
`
	blocks := certBlocks(raw)
	if len(blocks) != 2 {
		t.Fatalf("certBlocks() returned %d blocks, want 2", len(blocks))
	}
	for i, b := range blocks {
		if got := b[:27]; got != "-----BEGIN CERTIFICATE-----" {
			t.Errorf("block %d does not start with PEM header: %q", i, got)
		}
	}
}

func TestCertBlocksNone(t *testing.T) {
	if blocks := certBlocks("no certs here"); len(blocks) != 0 {
		t.Errorf("certBlocks() = %d blocks, want 0", len(blocks))
	}
}

func TestDefaultSaveName(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"api.example.com", "api.example.com_ca_chain.pem"},
		{"*.example.com", "example.com_ca_chain.pem"},
		{"", "cert_ca_chain.pem"},
		{"My CA Name", "My_CA_Name_ca_chain.pem"},
	}
	for _, tt := range tests {
		if got := defaultSaveName(tt.in); got != tt.want {
			t.Errorf("defaultSaveName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
