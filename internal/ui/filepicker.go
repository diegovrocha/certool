package ui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type FilePicker struct {
	Prompt   string
	allFiles []string
	filtered []string
	cursor   int
	filter   textinput.Model
	Selected string
	Done     bool
}

func NewFilePicker(prompt string, findArgs ...string) FilePicker {
	ti := textinput.New()
	ti.Placeholder = "type to filter..."
	ti.Focus()

	files := findFiles(findArgs...)

	return FilePicker{
		Prompt:   prompt,
		allFiles: files,
		filtered: files,
		filter:   ti,
	}
}

func NewCertFilePicker(prompt string) FilePicker {
	return NewFilePicker(prompt,
		"-iname", "*.pfx", "-o", "-iname", "*.p12",
		"-o", "-iname", "*.pem", "-o", "-iname", "*.cer",
		"-o", "-iname", "*.crt",
	)
}

func NewAllFilePicker(prompt string) FilePicker {
	return NewFilePicker(prompt,
		"-iname", "*.pfx", "-o", "-iname", "*.p12",
		"-o", "-iname", "*.pem", "-o", "-iname", "*.cer",
		"-o", "-iname", "*.crt", "-o", "-iname", "*.key",
	)
}

func NewPfxFilePicker(prompt string) FilePicker {
	return NewFilePicker(prompt, "-iname", "*.pfx", "-o", "-iname", "*.p12")
}

func NewCertOnlyPicker(prompt string) FilePicker {
	return NewFilePicker(prompt,
		"-iname", "*.pem", "-o", "-iname", "*.crt",
		"-o", "-iname", "*.cer",
	)
}

func NewKeyPicker(prompt string) FilePicker {
	return NewFilePicker(prompt, "-iname", "*.key", "-o", "-iname", "*.pem")
}

func (fp FilePicker) Init() tea.Cmd {
	return textinput.Blink
}

func (fp FilePicker) Update(msg tea.Msg) (FilePicker, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if fp.cursor > 0 {
				fp.cursor--
			}
			return fp, nil
		case "down":
			if fp.cursor < len(fp.filtered)-1 {
				fp.cursor++
			}
			return fp, nil
		case "enter":
			if len(fp.filtered) > 0 {
				fp.Selected = fp.filtered[fp.cursor]
				fp.Done = true
			}
			return fp, nil
		}
	}

	// Update the text filter
	var cmd tea.Cmd
	fp.filter, cmd = fp.filter.Update(msg)

	// Apply filter
	query := strings.ToLower(fp.filter.Value())
	if query == "" {
		fp.filtered = fp.allFiles
	} else {
		fp.filtered = nil
		for _, f := range fp.allFiles {
			if strings.Contains(strings.ToLower(f), query) {
				fp.filtered = append(fp.filtered, f)
			}
		}
	}

	// Adjust cursor
	if fp.cursor >= len(fp.filtered) {
		fp.cursor = len(fp.filtered) - 1
	}
	if fp.cursor < 0 {
		fp.cursor = 0
	}

	return fp, cmd
}

func (fp FilePicker) View() string {
	var b strings.Builder

	b.WriteString("  " + ActiveStyle.Render(fp.Prompt) + "\n\n")
	b.WriteString("  " + fp.filter.View() + "\n\n")

	if len(fp.allFiles) == 0 {
		b.WriteString("  " + ErrorStyle.Render("No files found") + "\n")
		return b.String()
	}

	if len(fp.filtered) == 0 {
		b.WriteString("  " + DimStyle.Render("No results for this filter") + "\n")
		return b.String()
	}

	// Show at most 15 items with scroll
	maxVisible := 15
	start := 0
	if fp.cursor >= maxVisible {
		start = fp.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(fp.filtered) {
		end = len(fp.filtered)
	}

	if start > 0 {
		b.WriteString(fmt.Sprintf("  %s\n", DimStyle.Render(fmt.Sprintf("  ↑ %d more above", start))))
	}

	for i := start; i < end; i++ {
		f := fp.filtered[i]
		if i == fp.cursor {
			b.WriteString(fmt.Sprintf("  %s%s\n", ActiveStyle.Render("➤ "), ActiveStyle.Render(f)))
		} else {
			b.WriteString(fmt.Sprintf("    %s\n", f))
		}
	}

	remaining := len(fp.filtered) - end
	if remaining > 0 {
		b.WriteString(fmt.Sprintf("  %s\n", DimStyle.Render(fmt.Sprintf("  ↓ %d more below", remaining))))
	}

	b.WriteString(fmt.Sprintf("\n  %s\n", DimStyle.Render(fmt.Sprintf("%d of %d files", len(fp.filtered), len(fp.allFiles)))))

	return b.String()
}

func findFiles(args ...string) []string {
	findArgs := []string{".", "-maxdepth", "5", "-type", "f", "("}
	findArgs = append(findArgs, args...)
	findArgs = append(findArgs, ")")
	out, err := exec.Command("find", findArgs...).Output()
	if err != nil {
		return nil
	}
	var result []string
	for _, l := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		l = strings.TrimPrefix(l, "./")
		if l != "" {
			result = append(result, l)
		}
	}
	return result
}
