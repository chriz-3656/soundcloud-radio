package lyrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LRCLibResponse struct {
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

// Fetch attempts to get lyrics for a given track artist and title.
func Fetch(artist, title string) (string, error) {
	// Clean up title (remove 'feat.', parentheses, etc. for better search accuracy)
	title = cleanString(title)
	artist = cleanString(artist)

	query := url.QueryEscape(fmt.Sprintf("%s %s", artist, title))
	apiURL := fmt.Sprintf("https://lrclib.net/api/search?q=%s", query)

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "SoundCloud-Radio-TUI/1.0 (https://github.com/chriz-3656/soundcloud-radio)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("lrclib returned status %d", resp.StatusCode)
	}

	var results []LRCLibResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no lyrics found")
	}

	// Prefer synced, fallback to plain
	best := results[0]
	if best.PlainLyrics != "" {
		return best.PlainLyrics, nil
	}
	if best.SyncedLyrics != "" {
		return best.SyncedLyrics, nil
	}

	return "", fmt.Errorf("empty lyrics block")
}

func cleanString(s string) string {
	s = strings.Split(s, "(feat.")[0]
	s = strings.Split(s, "(Feat.")[0]
	s = strings.Split(s, " ft.")[0]
	s = strings.Split(s, " Ft.")[0]
	s = strings.Split(s, "[")[0]
	return strings.TrimSpace(s)
}
