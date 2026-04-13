package menu

import (
	"strings"
	"testing"
)

func TestNewMenu(t *testing.T) {
	m := New()

	if len(m.items) == 0 {
		t.Error("Menu should have items")
	}

	// First selectable item should not be a separator
	if m.items[m.cursor].isSeparator {
		t.Error("Initial cursor should not be on a separator")
	}
}

func TestMoveCursorSkipsSeparators(t *testing.T) {
	m := New()

	// Cursor starts at 1 (first selectable item)
	initial := m.cursor
	if m.items[initial].isSeparator {
		t.Error("Initial position should not be a separator")
	}

	// Move down
	m.moveCursor(1)
	if m.items[m.cursor].isSeparator {
		t.Error("moveCursor(1) should not stop on a separator")
	}

	// Move up
	m.moveCursor(-1)
	if m.items[m.cursor].isSeparator {
		t.Error("moveCursor(-1) should not stop on a separator")
	}
}

func TestMoveCursorWraps(t *testing.T) {
	m := New()

	// Move up many times - should wrap
	for i := 0; i < len(m.items)*2; i++ {
		m.moveCursor(-1)
		if m.items[m.cursor].isSeparator {
			t.Errorf("Cursor stopped on separator at index %d", m.cursor)
		}
	}

	// Move down many times - should wrap
	for i := 0; i < len(m.items)*2; i++ {
		m.moveCursor(1)
		if m.items[m.cursor].isSeparator {
			t.Errorf("Cursor stopped on separator at index %d", m.cursor)
		}
	}
}

func TestMenuHasAllActions(t *testing.T) {
	m := New()

	expectedActions := []string{
		"pfx_pem", "pfx_cer_pem", "pfx_cer_der", "pfx_key", "pfx_repack",
		"inspect", "verify_chain", "verify_key", "compare_hash", "gen_self", "quit",
	}

	actions := make(map[string]bool)
	for _, item := range m.items {
		if item.action != "" {
			actions[item.action] = true
		}
	}

	for _, expected := range expectedActions {
		if !actions[expected] {
			t.Errorf("Menu should have action '%s'", expected)
		}
	}
}

func TestMenuView(t *testing.T) {
	m := New()
	v := m.View()

	if !strings.Contains(v, "CONVERT") {
		t.Error("View should contain CONVERT section")
	}
	if !strings.Contains(v, "VALIDATE") {
		t.Error("View should contain VALIDATE section")
	}
	if !strings.Contains(v, "GENERATE") {
		t.Error("View should contain GENERATE section")
	}
	if !strings.Contains(v, "➤") {
		t.Error("View should contain cursor ➤")
	}
	if !strings.Contains(v, "Quit") {
		t.Error("View should contain Quit option")
	}
}

func TestMenuQuit(t *testing.T) {
	m := New()
	m.quitting = true
	v := m.View()

	if !strings.Contains(v, "Goodbye") {
		t.Error("View on quit should show 'Goodbye'")
	}
}
