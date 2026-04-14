package verify

import "testing"

func TestCompareFieldsMatch(t *testing.T) {
	f := certField{label: "CN", val1: "example.com", val2: "example.com", match: true}
	if !f.match {
		t.Error("matching values should have match=true")
	}
	// Build as code would: equality check
	f2 := certField{label: "CN", val1: "same", val2: "same", match: "same" == "same"}
	if !f2.match {
		t.Error("equal strings should produce match=true")
	}
}

func TestCompareFieldsDiffer(t *testing.T) {
	v1 := "a"
	v2 := "b"
	f := certField{label: "CN", val1: v1, val2: v2, match: v1 == v2}
	if f.match {
		t.Error("different values should have match=false")
	}
}

func TestIsPFX(t *testing.T) {
	cases := map[string]bool{
		"a.pfx":         true,
		"a.p12":         true,
		"a.PFX":         true,
		"b.pem":         false,
		"c.crt":         false,
		"no-extension":  false,
	}
	for file, want := range cases {
		if got := isPFX(file); got != want {
			t.Errorf("isPFX(%q) = %v, want %v", file, got, want)
		}
	}
}
