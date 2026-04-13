package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func OpenSSLVersion() string {
	out, err := exec.Command("openssl", "version").Output()
	if err != nil {
		return "OpenSSL not found"
	}
	return strings.TrimSpace(string(out))
}

func ResultBox(success bool, title string, lines ...string) string {
	var style lipgloss.Style
	var icon string

	if success {
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGreen).
			Padding(0, 2)
		icon = SuccessStyle.Render("✔ " + title)
	} else {
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorRed).
			Padding(0, 2)
		icon = ErrorStyle.Render("✖ " + title)
	}

	content := icon + "\n"
	for _, line := range lines {
		content += "\n  " + line
	}

	return style.Render(content)
}

func CertBox(width int, lines ...string) string {
	inner := width - 6
	if inner < 40 {
		inner = 40
	}

	var b strings.Builder
	border := strings.Repeat("═", inner)

	b.WriteString(fmt.Sprintf("  ╔%s╗\n", border))
	for _, line := range lines {
		pad := inner - len(line)
		if pad < 0 {
			pad = 0
		}
		b.WriteString(fmt.Sprintf("  ║%s%*s║\n", line, pad, ""))
	}
	b.WriteString(fmt.Sprintf("  ╚%s╝\n", border))

	return b.String()
}
