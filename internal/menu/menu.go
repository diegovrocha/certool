package menu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diegovrocha/certui/internal/convert"
	"github.com/diegovrocha/certui/internal/generate"
	"github.com/diegovrocha/certui/internal/inspect"
	"github.com/diegovrocha/certui/internal/ui"
	"github.com/diegovrocha/certui/internal/update"
	"github.com/diegovrocha/certui/internal/verify"
)

type menuItem struct {
	label       string
	desc        string
	action      string
	isSeparator bool
}

type screen int

const (
	screenMenu screen = iota
	screenSub
)

type Model struct {
	items      []menuItem
	cursor     int
	screen     screen
	sub        tea.Model
	width      int
	height     int
	quitting   bool
	updateMsg  string
	updateDone bool
}

func New() Model {
	items := []menuItem{
		{label: "── CONVERT ──────────────────────────────────────", isSeparator: true},
		{label: "PFX/P12 → PEM", desc: "certificate + key as text", action: "pfx_pem"},
		{label: "PFX/P12 → CER", desc: "certificate PEM (text)", action: "pfx_cer_pem"},
		{label: "PFX/P12 → CER", desc: "certificate DER (binary)", action: "pfx_cer_der"},
		{label: "PFX/P12 → KEY", desc: "private key only", action: "pfx_key"},
		{label: "PFX/P12 → P12", desc: "repack --legacy → modern", action: "pfx_repack"},
		{label: "── VALIDATE ─────────────────────────────────────", isSeparator: true},
		{label: "Inspect", desc: "subject, validity, issuer...", action: "inspect"},
		{label: "Verify chain", desc: "validate cert → CA → root", action: "verify_chain"},
		{label: "Verify cert+key", desc: "check if cert matches key", action: "verify_key"},
		{label: "Compare certs", desc: "check if two certs are the same", action: "compare_hash"},
		{label: "── GENERATE ─────────────────────────────────────", isSeparator: true},
		{label: "Generate self-signed", desc: "create cert + key for dev/testing", action: "gen_self"},
		{label: "─────────────────────────────────────────────────", isSeparator: true},
		{label: "Update", desc: "download and install the latest version", action: "update"},
		{label: "Quit", action: "quit"},
	}

	m := Model{items: items, cursor: 1}
	return m
}

type updateCheckMsg string

func checkForUpdate() tea.Msg {
	return updateCheckMsg(ui.CheckUpdate())
}

func (m Model) Init() tea.Cmd {
	return checkForUpdate
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case updateCheckMsg:
		m.updateMsg = string(msg)
		m.updateDone = true
		return m, nil

	case tea.KeyMsg:
		if m.screen == screenSub {
			return m.updateSub(msg)
		}
		return m.updateMenu(msg)
	}

	if m.screen == screenSub && m.sub != nil {
		var cmd tea.Cmd
		m.sub, cmd = m.sub.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "enter":
		action := m.items[m.cursor].action
		return m.handleAction(action)
	}
	return m, nil
}

func (m Model) updateSub(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenMenu
		m.sub = nil
		return m, nil
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	}
	if m.sub != nil {
		var cmd tea.Cmd
		m.sub, cmd = m.sub.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) moveCursor(delta int) {
	n := len(m.items)
	for i := 0; i < n; i++ {
		m.cursor = (m.cursor + delta + n) % n
		if !m.items[m.cursor].isSeparator {
			break
		}
	}
}

func (m Model) handleAction(action string) (tea.Model, tea.Cmd) {
	switch action {
	case "quit":
		m.quitting = true
		return m, tea.Quit
	case "pfx_pem":
		m.screen = screenSub
		m.sub = convert.NewPfxToPem()
		return m, m.sub.Init()
	case "pfx_cer_pem":
		m.screen = screenSub
		m.sub = convert.NewPfxToCerPem()
		return m, m.sub.Init()
	case "pfx_cer_der":
		m.screen = screenSub
		m.sub = convert.NewPfxToCerDer()
		return m, m.sub.Init()
	case "pfx_key":
		m.screen = screenSub
		m.sub = convert.NewPfxToKey()
		return m, m.sub.Init()
	case "pfx_repack":
		m.screen = screenSub
		m.sub = convert.NewPfxRepack()
		return m, m.sub.Init()
	case "inspect":
		m.screen = screenSub
		m.sub = inspect.New()
		return m, m.sub.Init()
	case "verify_chain":
		m.screen = screenSub
		m.sub = verify.NewChain()
		return m, m.sub.Init()
	case "verify_key":
		m.screen = screenSub
		m.sub = verify.NewCertKey()
		return m, m.sub.Init()
	case "compare_hash":
		m.screen = screenSub
		m.sub = verify.NewCompareHash()
		return m, m.sub.Init()
	case "gen_self":
		m.screen = screenSub
		m.sub = generate.NewSelfSigned()
		return m, m.sub.Init()
	case "update":
		m.screen = screenSub
		m.sub = update.New()
		return m, m.sub.Init()
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "\n  " + ui.SuccessStyle.Render("Goodbye!") + "\n\n"
	}

	if m.screen == screenSub && m.sub != nil {
		return m.sub.View()
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(ui.Banner())
	b.WriteString("  " + ui.DimStyle.Render(ui.OpenSSLVersion()) + "\n")
	if m.updateMsg != "" {
		b.WriteString("  " + ui.WarnStyle.Render(m.updateMsg) + "\n")
	}
	b.WriteString("\n")

	// Build menu lines
	for _, item := range m.items {
		if item.isSeparator {
			b.WriteString(fmt.Sprintf("  %s\n", ui.SeparatorStyle.Render(item.label)))
			continue
		}

		cursor := "  "
		labelStyle := ui.InactiveStyle
		if m.items[m.cursor].action == item.action && m.cursor == m.indexOf(item) {
			cursor = ui.ActiveStyle.Render("➤ ")
			labelStyle = ui.ActiveStyle
		}

		label := labelStyle.Render(fmt.Sprintf("%-20s", item.label))
		desc := ui.DescStyle.Render(item.desc)
		b.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, label, desc))
	}

	b.WriteString("\n  " + ui.DimStyle.Render("↑/↓ navigate  enter select  q / ctrl+c quit") + "\n")

	return b.String()
}

func (m Model) indexOf(target menuItem) int {
	for i, item := range m.items {
		if item.action == target.action && item.label == target.label {
			return i
		}
	}
	return -1
}
