package ui

import (
	"fmt"
	"os/exec"
	"strings"
)

type FileStats struct {
	Total int
	PFX   int
	PEM   int
	CER   int
	KEY   int
}

func CountCertFiles() FileStats {
	var s FileStats

	count := func(args ...string) int {
		findArgs := []string{".", "-maxdepth", "5", "-type", "f", "("}
		findArgs = append(findArgs, args...)
		findArgs = append(findArgs, ")")
		out, err := exec.Command("find", findArgs...).Output()
		if err != nil {
			return 0
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			return 0
		}
		return len(lines)
	}

	s.PFX = count("-iname", "*.pfx", "-o", "-iname", "*.p12")
	s.PEM = count("-iname", "*.pem", "-o", "-iname", "*.crt")
	s.CER = count("-iname", "*.cer", "-o", "-iname", "*.der")
	s.KEY = count("-iname", "*.key")
	s.Total = s.PFX + s.PEM + s.CER + s.KEY

	return s
}

func FormatSidebar(stats FileStats) []string {
	return []string{
		"┌─ Quick Info ──────────┐",
		fmt.Sprintf("│ Files found: %-8d │", stats.Total),
		fmt.Sprintf("│ .pfx/.p12:   %-8d │", stats.PFX),
		fmt.Sprintf("│ .pem/.crt:   %-8d │", stats.PEM),
		fmt.Sprintf("│ .cer/.der:   %-8d │", stats.CER),
		fmt.Sprintf("│ .key:        %-8d │", stats.KEY),
		"└───────────────────────┘",
	}
}
