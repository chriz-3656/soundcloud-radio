package ui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type artworkMsg struct {
	trackID int64
	art     string
}

func fetchArtworkCmd(trackID int64, url string) tea.Cmd {
	return func() tea.Msg {
		if url == "" {
			return artworkMsg{trackID: trackID, art: ""}
		}

		// SoundCloud returns -large (100x100). We can request -t200x200 or just use large.
		url = strings.Replace(url, "-large.jpg", "-t200x200.jpg", 1)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return artworkMsg{trackID: trackID, art: ""}
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			return artworkMsg{trackID: trackID, art: ""}
		}

		bounds := img.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		if w == 0 || h == 0 {
			return artworkMsg{trackID: trackID, art: ""}
		}

		targetW, targetH := 28, 26 // 28 chars wide, 13 chars tall (since each char is 2 vertical pixels)
		
		var sb strings.Builder
		for y := 0; y < targetH; y += 2 {
			for x := 0; x < targetW; x++ {
				// Top pixel
				srcX := x * w / targetW
				srcY1 := y * h / targetH
				r1, g1, b1, _ := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY1).RGBA()
				
				// Bottom pixel
				srcY2 := (y + 1) * h / targetH
				r2, g2, b2, _ := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY2).RGBA()
				
				r1, g1, b1 = r1>>8, g1>>8, b1>>8
				r2, g2, b2 = r2>>8, g2>>8, b2>>8
				
				sb.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", r1, g1, b1, r2, g2, b2))
			}
			sb.WriteString("\x1b[0m\n") // reset color at end of line
		}
		
		return artworkMsg{trackID: trackID, art: sb.String()}
	}
}
