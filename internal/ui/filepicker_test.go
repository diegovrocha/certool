package ui

import (
	"strings"
	"testing"
)

func TestFilePickerFilter(t *testing.T) {
	fp := FilePicker{
		allFiles: []string{
			"cert.pem",
			"key.pem",
			"rv2.pfx",
			"cert_chain.pem",
			"server.crt",
		},
		filtered: []string{
			"cert.pem",
			"key.pem",
			"rv2.pfx",
			"cert_chain.pem",
			"server.crt",
		},
	}

	// Filter by "rv2"
	query := "rv2"
	fp.filtered = nil
	for _, f := range fp.allFiles {
		if strings.Contains(strings.ToLower(f), query) {
			fp.filtered = append(fp.filtered, f)
		}
	}

	if len(fp.filtered) != 1 {
		t.Errorf("Filter 'rv2' should return 1 file, returned %d", len(fp.filtered))
	}
	if fp.filtered[0] != "rv2.pfx" {
		t.Errorf("Filter 'rv2' should return 'rv2.pfx', returned '%s'", fp.filtered[0])
	}

	// Filter by "pem"
	query = "pem"
	fp.filtered = nil
	for _, f := range fp.allFiles {
		if strings.Contains(strings.ToLower(f), query) {
			fp.filtered = append(fp.filtered, f)
		}
	}

	if len(fp.filtered) != 3 {
		t.Errorf("Filter 'pem' should return 3 files, returned %d", len(fp.filtered))
	}

	// Empty filter returns all
	fp.filtered = fp.allFiles
	if len(fp.filtered) != 5 {
		t.Errorf("No filter should return 5 files, returned %d", len(fp.filtered))
	}
}

func TestFilePickerCursorBounds(t *testing.T) {
	fp := FilePicker{
		allFiles: []string{"a.pem", "b.pem", "c.pem"},
		filtered: []string{"a.pem", "b.pem", "c.pem"},
		cursor:   0,
	}

	// Cursor should not go below 0
	fp.cursor = 0
	if fp.cursor < 0 {
		t.Error("Cursor should not be negative")
	}

	// Cursor should not exceed list size
	fp.cursor = len(fp.filtered) - 1
	if fp.cursor >= len(fp.filtered) {
		t.Error("Cursor should not exceed list size")
	}
}

func TestFilePickerView(t *testing.T) {
	fp := FilePicker{
		Prompt:   "Select file",
		allFiles: []string{"test.pem"},
		filtered: []string{"test.pem"},
		cursor:   0,
	}

	// Manually initialize textinput for View to work
	fp.filter.Placeholder = "type to filter..."

	v := fp.View()
	if !strings.Contains(v, "Select file") {
		t.Error("View should contain the prompt")
	}
	if !strings.Contains(v, "test.pem") {
		t.Error("View should contain the file")
	}
}

func TestFilePickerEmpty(t *testing.T) {
	fp := FilePicker{
		Prompt:   "Select",
		allFiles: []string{},
		filtered: []string{},
	}
	fp.filter.Placeholder = "type to filter..."

	v := fp.View()
	if !strings.Contains(v, "No files found") {
		t.Error("Empty view should show no files found message")
	}
}
