package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	repoOwner = "diegovrocha"
	repoName  = "certool"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

// CheckUpdate queries the GitHub API for the latest release.
// Returns the update message to display, or empty string if up to date.
// Never blocks longer than 2 seconds; errors are silently ignored.
func CheckUpdate() string {
	ch := make(chan string, 1)

	go func() {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			ch <- ""
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			ch <- ""
			return
		}

		var release githubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			ch <- ""
			return
		}

		latest := strings.TrimPrefix(release.TagName, "v")
		if latest != "" && latest != Version {
			ch <- fmt.Sprintf("  Update v%s available! Run: brew upgrade certool", latest)
		} else {
			ch <- ""
		}
	}()

	// Wait up to 2 seconds for the result
	select {
	case msg := <-ch:
		return msg
	case <-time.After(2 * time.Second):
		return ""
	}
}
