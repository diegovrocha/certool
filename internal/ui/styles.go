package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	ColorCyan    = lipgloss.Color("14")
	ColorMagenta = lipgloss.Color("5")
	ColorGreen   = lipgloss.Color("2")
	ColorRed     = lipgloss.Color("1")
	ColorYellow  = lipgloss.Color("3")
	ColorDim     = lipgloss.Color("8")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan)

	SubtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(ColorMagenta)

	DimStyle = lipgloss.NewStyle().
			Faint(true)

	ActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan)

	InactiveStyle = lipgloss.NewStyle()

	SeparatorStyle = lipgloss.NewStyle().
			Faint(true)

	DescStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGreen)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorRed)

	WarnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorYellow)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorCyan).
			Padding(1, 2)
)

// Version is set at build time via ldflags by GoReleaser.
// Falls back to "dev" when built locally without ldflags.
var Version = "dev"

func Banner() string {
	t := TitleStyle.Render
	d := DimStyle.Render
	s := SubtitleStyle.Render

	// All logo lines padded to 31 chars so right-side text aligns
	var b strings.Builder
	b.WriteString(t("  ____         _____ _   _ ___ ") + "\n")
	b.WriteString(t(" / ___|___ _ _|_   _| | | |_ _|") + "  " + s("Cert + TUI") + "\n")
	b.WriteString(t("| |   / _ \\ '__|| | | | | || | ") + "  " + s("Digital certificate conversion,") + "\n")
	b.WriteString(t("| |__|  __/ |   | | | |_| || | ") + "  " + s("validation and generation.") + "\n")
	b.WriteString(t(" \\____\\___|_|   |_|  \\___/|___|") + "  " + d("https://github.com/diegovrocha/certui") + "\n")
	b.WriteString(d("                                 v"+Version) + "\n")
	return b.String()
}
