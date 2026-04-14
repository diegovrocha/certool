package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diegovrocha/certui/internal/ui"
)

type step int

const (
	stepConfirm step = iota
	stepDownloading
	stepDone
)

type Model struct {
	step      step
	current   string
	latest    string
	body      string
	scroll    int
	result    string
	success   bool
}

type updateInfoMsg struct {
	latest string
	body   string
	err    string
}

type downloadResultMsg struct {
	success bool
	message string
}

func New() tea.Model {
	return &Model{
		step:    stepConfirm,
		current: ui.Version,
	}
}

func (m *Model) Init() tea.Cmd {
	return fetchLatestVersion
}

func fetchLatestVersion() tea.Msg {
	url := "https://api.github.com/repos/diegovrocha/certui/releases/latest"
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return updateInfoMsg{err: "Could not reach GitHub: " + err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return updateInfoMsg{err: fmt.Sprintf("GitHub API returned %d", resp.StatusCode)}
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return updateInfoMsg{err: "Invalid response from GitHub"}
	}

	return updateInfoMsg{latest: strings.TrimPrefix(release.TagName, "v"), body: release.Body}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateInfoMsg:
		if msg.err != "" {
			m.result = msg.err
			m.success = false
			m.step = stepDone
			return m, nil
		}
		m.latest = msg.latest
		m.body = msg.body
		return m, nil

	case downloadResultMsg:
		m.success = msg.success
		m.result = msg.message
		m.step = stepDone
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, nil
		}
		if m.step == stepConfirm {
			switch msg.String() {
			case "up", "k":
				if m.scroll > 0 {
					m.scroll--
				}
				return m, nil
			case "down", "j":
				m.scroll++
				return m, nil
			}
		}
		if m.step == stepConfirm && msg.String() == "enter" {
			if m.latest == "" {
				return m, nil
			}
			if normalizeVer(m.latest) == normalizeVer(m.current) {
				m.success = true
				m.result = "Already on the latest version"
				m.step = stepDone
				return m, nil
			}
			m.step = stepDownloading
			return m, m.doUpdate()
		}
	}
	return m, nil
}

func (m *Model) doUpdate() tea.Cmd {
	return func() tea.Msg {
		os := runtime.GOOS
		arch := runtime.GOARCH

		ext := "tar.gz"
		if os == "windows" {
			ext = "zip"
		}

		url := fmt.Sprintf("https://github.com/diegovrocha/certui/releases/latest/download/certui_%s_%s.%s", os, arch, ext)

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return downloadResultMsg{false, "Download failed: " + err.Error()}
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return downloadResultMsg{false, fmt.Sprintf("Download failed: HTTP %d", resp.StatusCode)}
		}

		// Current executable path
		exe, err := getExecutable()
		if err != nil {
			return downloadResultMsg{false, "Could not locate current binary: " + err.Error()}
		}

		// Write download to temp
		tmpDir, err := osMkdirTemp("", "certui-update-*")
		if err != nil {
			return downloadResultMsg{false, "Could not create temp dir: " + err.Error()}
		}
		defer removeAll(tmpDir)

		archivePath := filepath.Join(tmpDir, "archive."+ext)
		out, err := createFile(archivePath)
		if err != nil {
			return downloadResultMsg{false, "Could not save download: " + err.Error()}
		}
		if _, err := io.Copy(out, resp.Body); err != nil {
			out.Close()
			return downloadResultMsg{false, "Download write failed: " + err.Error()}
		}
		out.Close()

		// Extract certui binary from archive
		binName := "certui"
		if os == "windows" {
			binName = "certui.exe"
		}
		extracted := filepath.Join(tmpDir, binName)

		if ext == "zip" {
			if err := extractZip(archivePath, binName, extracted); err != nil {
				return downloadResultMsg{false, "Extract failed: " + err.Error()}
			}
		} else {
			if err := extractTarGz(archivePath, binName, extracted); err != nil {
				return downloadResultMsg{false, "Extract failed: " + err.Error()}
			}
		}

		// Replace current binary
		// On Unix we can overwrite a running binary; on Windows we need a workaround
		if err := replaceBinary(extracted, exe); err != nil {
			return downloadResultMsg{false, "Replace failed: " + err.Error()}
		}

		return downloadResultMsg{true, fmt.Sprintf("Updated to v%s. Restart certui to use the new version.", m.latest)}
	}
}

func extractTarGz(archive, wantName, dest string) error {
	f, err := openFile(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return fmt.Errorf("binary %s not found in archive", wantName)
		}
		if err != nil {
			return err
		}
		if filepath.Base(h.Name) == wantName {
			out, err := createFile(dest)
			if err != nil {
				return err
			}
			defer out.Close()
			if _, err := io.Copy(out, tr); err != nil {
				return err
			}
			return chmod(dest, 0755)
		}
	}
}

func extractZip(archive, wantName, dest string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == wantName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			out, err := createFile(dest)
			if err != nil {
				return err
			}
			defer out.Close()
			if _, err := io.Copy(out, rc); err != nil {
				return err
			}
			return chmod(dest, 0755)
		}
	}
	return fmt.Errorf("binary %s not found in archive", wantName)
}

func replaceBinary(src, dst string) error {
	// Resolve symlinks so we overwrite the real binary
	real, err := evalSymlinks(dst)
	if err == nil {
		dst = real
	}

	// Try direct rename first (atomic on same filesystem)
	if err := rename(src, dst); err == nil {
		return nil
	}

	// Fallback: copy content (handles cross-device)
	in, err := openFile(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// On the same path, we need to write while the running binary may have it mapped.
	// On Unix, unlink first then recreate works.
	_ = removeFile(dst)
	out, err := createFile(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return chmod(dst, 0755)
}

func normalizeVer(v string) string {
	v = strings.TrimRight(v, ".0")
	if v == "" {
		v = "0"
	}
	return v
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(ui.Banner())
	b.WriteString("\n  " + ui.TitleStyle.Render("── Update ──") + "\n\n")

	b.WriteString(fmt.Sprintf("  Current version:  %s\n", ui.ActiveStyle.Render("v"+m.current)))

	switch m.step {
	case stepConfirm:
		if m.latest == "" {
			b.WriteString("  Latest version:   " + ui.DimStyle.Render("checking...") + "\n")
		} else {
			b.WriteString(fmt.Sprintf("  Latest version:   %s\n", ui.ActiveStyle.Render("v"+m.latest)))
			if normalizeVer(m.latest) == normalizeVer(m.current) {
				b.WriteString("\n  " + ui.SuccessStyle.Render("✔ You are on the latest version") + "\n")
			} else {
				b.WriteString("\n  " + ui.WarnStyle.Render("⚠ Update available") + "\n")
				// Changelog
				if strings.TrimSpace(m.body) != "" {
					b.WriteString("\n  " + ui.TitleStyle.Render("Changelog:") + "\n")
					lines := strings.Split(strings.ReplaceAll(m.body, "\r\n", "\n"), "\n")
					const maxLines = 15
					total := len(lines)
					if m.scroll < 0 {
						m.scroll = 0
					}
					if total <= maxLines {
						m.scroll = 0
						for _, l := range lines {
							b.WriteString("    " + ui.DimStyle.Render(l) + "\n")
						}
					} else {
						maxScroll := total - maxLines
						if m.scroll > maxScroll {
							m.scroll = maxScroll
						}
						end := m.scroll + maxLines
						if end > total {
							end = total
						}
						if m.scroll > 0 {
							b.WriteString("    " + ui.DimStyle.Render(fmt.Sprintf("↑ %d lines above", m.scroll)) + "\n")
						}
						for i := m.scroll; i < end; i++ {
							b.WriteString("    " + ui.DimStyle.Render(lines[i]) + "\n")
						}
						remaining := total - end
						if remaining > 0 {
							b.WriteString("    " + ui.DimStyle.Render(fmt.Sprintf("↓ %d lines below (↑/↓ to scroll)", remaining)) + "\n")
						}
					}
				}
				b.WriteString("\n  " + ui.DimStyle.Render("Press ENTER to install, esc to cancel") + "\n")
			}
		}

	case stepDownloading:
		b.WriteString(fmt.Sprintf("  Latest version:   %s\n\n", ui.ActiveStyle.Render("v"+m.latest)))
		b.WriteString("  ⏳ Downloading and installing...\n")

	case stepDone:
		b.WriteString("\n")
		if m.success {
			b.WriteString(ui.ResultBox(true, "Success", m.result))
		} else {
			b.WriteString(ui.ResultBox(false, "Error", m.result))
		}
	}

	b.WriteString("\n  " + ui.DimStyle.Render("esc back  enter confirm  ctrl+c quit") + "\n")
	return b.String()
}

// helpers to keep imports clean

func getExecutable() (string, error)        { return osExecutable() }
func osExecutable() (string, error)         { return osExecutableReal() }
func osMkdirTemp(d, p string) (string, error) { return os.MkdirTemp(d, p) }
func removeAll(p string)                    { _ = os.RemoveAll(p) }
func createFile(p string) (*os.File, error) { return os.Create(p) }
func openFile(p string) (*os.File, error)   { return os.Open(p) }
func chmod(p string, m os.FileMode) error   { return os.Chmod(p, m) }
func rename(src, dst string) error          { return os.Rename(src, dst) }
func removeFile(p string) error             { return os.Remove(p) }
func evalSymlinks(p string) (string, error) { return filepath.EvalSymlinks(p) }

func osExecutableReal() (string, error) {
	return os.Executable()
}
