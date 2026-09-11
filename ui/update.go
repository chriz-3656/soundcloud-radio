package ui

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"soundcloud-radio/internal/lyrics"
	"soundcloud-radio/internal/player"
	"soundcloud-radio/internal/resolver"
	"soundcloud-radio/internal/soundcloud"
)

type searchResultMsg []soundcloud.Track
type errMsg error
type streamResolvedMsg struct {
	track  soundcloud.Track
	stream *resolver.ResolvedStream
}
type relatedTracksMsg []soundcloud.Track

type setupLogMsg string
type setupDoneMsg struct{}
type setupStepMsg int

type lyricsMsg struct {
	trackID int64
	lyrics  string
}

type lyricsCheckResultMsg struct {
	trackID   int64
	hasLyrics bool
}

func (m *Model) fetchLyricsCmd(artist, title string, trackID int64) tea.Cmd {
	return func() tea.Msg {
		l, err := lyrics.Fetch(artist, title)
		if err != nil {
			return lyricsMsg{trackID: trackID, lyrics: ""}
		}
		return lyricsMsg{trackID: trackID, lyrics: l}
	}
}

func (m *Model) checkLyricsBulkCmd(tracks []soundcloud.Track) tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range tracks {
		if !m.lyricsChecked[t.ID] {
			track := t
			cmds = append(cmds, func() tea.Msg {
				l, err := lyrics.Fetch(track.Artist, track.Title)
				hasLyrics := (err == nil && l != "")
				return lyricsCheckResultMsg{trackID: track.ID, hasLyrics: hasLyrics}
			})
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) showNotif(text string) {
	m.notification = text
	m.notifTimer = 12 // 3 seconds at 4 ticks/sec
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searchInput.Focused() && msg.String() != "enter" && msg.String() != "esc" && msg.String() != "ctrl+c" {
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			m.cancel()
			m.player.Close()
			return m, tea.Quit
		case "esc":
			if m.searchInput.Focused() {
				m.searchInput.Blur()
			} else {
				m.viewState = ViewHome
			}
		case "ctrl+f", "/":
			m.viewState = ViewSearch
			m.searchInput.Focus()
		case " ":
			m.player.TogglePause()
		case "n":
			if track, ok := m.queue.Pop(); ok {
				m.statusMsg = "Resolving next track..."
				cmds = append(cmds, m.resolveCmd(track))
			}
		case "p":
			hist := m.store.GetHistory()
			if len(hist) > 0 {
				prev := hist[0].Track
				if m.currentTrack != nil {
					m.queue.PlayNext(*m.currentTrack)
				}
				m.statusMsg = "Resolving " + prev.Title + "..."
				cmds = append(cmds, m.resolveCmd(prev))
			}
		case "s":
			m.player.Stop()
			m.currentTrack = nil
		case "q":
			m.cancel()
			m.player.Close()
			return m, tea.Quit
		case "+", "=":
			m.player.SetVolume(m.player.Volume() + 5)
			m.showNotif(fmt.Sprintf("Volume: %d%%", m.player.Volume()))
		case "-", "_":
			m.player.SetVolume(m.player.Volume() - 5)
			m.showNotif(fmt.Sprintf("Volume: %d%%", m.player.Volume()))
		case "m":
			if m.player.Volume() > 0 {
				m.player.SetVolume(0)
				m.showNotif("Muted")
			} else {
				m.player.SetVolume(100)
				m.showNotif("Unmuted")
			}
		case "l", "L":
			if m.hasLyrics {
				if m.viewState == ViewLyrics {
					m.viewState = ViewHome
				} else {
					m.viewState = ViewLyrics
				}
			} else {
				m.showNotif("No lyrics available for this track")
			}
		case "f", "F":
			if m.currentTrack != nil {
				if m.store.IsFavorite(m.currentTrack.ID) {
					m.store.RemoveFavorite(m.currentTrack.ID)
					m.showNotif("✓ Removed from favorites")
				} else {
					m.store.AddFavorite(*m.currentTrack)
					m.showNotif("✓ Added to favorites")
				}
			}
		case "?":
			m.viewState = ViewHelp
		}

		// View specific navigation
		switch m.viewState {
		case ViewSplash:
			if msg.String() == "enter" {
				m.viewState = ViewSetup
				m.setupLogs = []string{"[*] Initializing environment checks..."}
				return m, func() tea.Msg { return setupStepMsg(1) }
			}
		case ViewSearch:
			switch msg.String() {
			case "enter":
				if m.searchInput.Focused() {
					query := m.searchInput.Value()
					if query != "" {
						m.isSearching = true
						m.searchInput.Blur()
						cmds = append(cmds, m.searchCmd(query))
					}
				} else {
					if len(m.searchResults) > 0 {
						selected := m.searchResults[m.searchCursor]
						m.queue.PlayNext(selected)
						if track, ok := m.queue.Pop(); ok {
							m.statusMsg = "Resolving " + track.Title + "..."
							cmds = append(cmds, m.resolveCmd(track))
						}
					}
				}
			case "a":
				if len(m.searchResults) > 0 {
					m.queue.Add(m.searchResults[m.searchCursor])
					m.showNotif("✓ Added to queue")
				}
			case "r":
				m.cfg.Radio = true
				m.showNotif("✓ Radio enabled")
			case "j", "down":
				if m.searchCursor < len(m.searchResults)-1 {
					m.searchCursor++
				}
			case "k", "up":
				if m.searchCursor > 0 {
					m.searchCursor--
				}
			}

		case ViewHome: // Used as Queue view
			switch msg.String() {
			case "j", "down":
				if m.queueCursor < m.queue.Len()-1 {
					m.queueCursor++
				}
			case "k", "up":
				if m.queueCursor > 0 {
					m.queueCursor--
				}
			case "d":
				if m.queue.Len() > 0 {
					m.queue.Remove(m.queueCursor)
					if m.queueCursor >= m.queue.Len() && m.queueCursor > 0 {
						m.queueCursor--
					}
					m.showNotif("✓ Removed from queue")
				}
			case "c":
				m.queue.Clear()
				m.queueCursor = 0
				m.showNotif("✓ Queue cleared")
			case "S": // Shift+S for shuffle
				m.queue.Shuffle()
				m.showNotif("✓ Queue shuffled")
			case "enter":
				// Play selected from queue
				if m.queue.Len() > 0 {
					items := m.queue.Items()
					selected := items[m.queueCursor]
					m.queue.Remove(m.queueCursor)
					m.queue.PlayNext(selected)
					if track, ok := m.queue.Pop(); ok {
						m.statusMsg = "Resolving " + track.Title + "..."
						cmds = append(cmds, m.resolveCmd(track))
					}
				}
			}

		case ViewFavorites:
			favs := m.store.GetFavorites()
			switch msg.String() {
			case "j", "down":
				if m.favCursor < len(favs)-1 {
					m.favCursor++
				}
			case "k", "up":
				if m.favCursor > 0 {
					m.favCursor--
				}
			case "enter":
				if len(favs) > 0 {
					m.queue.PlayNext(favs[m.favCursor])
					if track, ok := m.queue.Pop(); ok {
						m.statusMsg = "Resolving " + track.Title + "..."
						cmds = append(cmds, m.resolveCmd(track))
					}
				}
			case "a":
				if len(favs) > 0 {
					m.queue.Add(favs[m.favCursor])
					m.showNotif("✓ Added to queue")
				}
			case "d":
				if len(favs) > 0 {
					m.store.RemoveFavorite(favs[m.favCursor].ID)
					if m.favCursor > 0 { m.favCursor-- }
					m.showNotif("✓ Removed from favorites")
				}
			}

		case ViewHistory:
			hist := m.store.GetHistory()
			switch msg.String() {
			case "j", "down":
				if m.histCursor < len(hist)-1 {
					m.histCursor++
				}
			case "k", "up":
				if m.histCursor > 0 {
					m.histCursor--
				}
			case "enter":
				if len(hist) > 0 {
					m.queue.PlayNext(hist[m.histCursor].Track)
					if track, ok := m.queue.Pop(); ok {
						m.statusMsg = "Resolving " + track.Title + "..."
						cmds = append(cmds, m.resolveCmd(track))
					}
				}
			case "a":
				if len(hist) > 0 {
					m.queue.Add(hist[m.histCursor].Track)
					m.showNotif("✓ Added to queue")
				}
			}

		case ViewLyrics:
			lines := strings.Split(m.currentLyrics, "\n")
			switch msg.String() {
			case "j", "down":
				if m.lyricsCursor < len(lines)-1 {
					m.lyricsCursor++
				}
			case "k", "up":
				if m.lyricsCursor > 0 {
					m.lyricsCursor--
				}
			}
		}

	case tea.MouseClickMsg:
		if true {
			if msg.X < 25 {
				switch msg.Y {
				case 3, 4: m.viewState = ViewHome
				case 5, 6: m.viewState = ViewSearch
				case 7, 8: m.viewState = ViewFavorites
				case 9, 10: m.viewState = ViewHistory
				case 12, 13: m.viewState = ViewHelp
				case 14, 15: m.viewState = ViewSettings
				}
			} else if msg.X >= 25 && msg.X < m.width-35 {
				if m.viewState == ViewSearch {
					if msg.Y >= 2 && msg.Y <= 4 {
						if !m.searchInput.Focused() {
							m.searchInput.Focus()
						}
					} else if !m.searchInput.Focused() && len(m.searchResults) > 0 {
						// Extremely approximate
						itemsPerPage := (m.height - 4 - 2 - 10) / 2
						if itemsPerPage < 1 { itemsPerPage = 1 }
						startIndex := 0
						if m.searchCursor >= itemsPerPage {
							startIndex = m.searchCursor - itemsPerPage + 1
						}
						
						idx := (msg.Y - 5) / 2
						if idx >= 0 && startIndex+idx < len(m.searchResults) {
							m.searchCursor = startIndex + idx
						}
					}
				}
			}
		}

	case setupStepMsg:
		if msg == 1 {
			return m, func() tea.Msg {
				_, err := exec.LookPath("mpv")
				if err == nil {
					return setupLogMsg("[✓] mpv audio backend found")
				}
				_, err = exec.LookPath("vlc")
				if err == nil {
					return setupLogMsg("[✓] vlc audio backend found (fallback)")
				}
				return setupLogMsg("[!] No audio backend found (Please install mpv for the best experience!)")
			}
		} else if msg == 2 {
			m.setupLogs = append(m.setupLogs, "[*] Checking yt-dlp stream resolver...")
			return m, tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg { return setupStepMsg(3) })
		} else if msg == 3 {
			return m, func() tea.Msg {
				_, err := resolver.EnsureYTDLP(m.ctx)
				if err != nil {
					return setupLogMsg("[✕] Failed to install yt-dlp: " + err.Error())
				}
				return setupLogMsg("[✓] yt-dlp resolver ready")
			}
		} else if msg == 4 {
			m.setupLogs = append(m.setupLogs, "\nAll checks completed. Booting UI...")
			return m, tea.Tick(time.Second*1, func(time.Time) tea.Msg { return setupDoneMsg{} })
		}

	case setupLogMsg:
		m.setupLogs = append(m.setupLogs, string(msg))
		if strings.Contains(string(msg), "audio backend") {
			return m, func() tea.Msg { return setupStepMsg(2) }
		}
		if strings.Contains(string(msg), "yt-dlp resolver") || strings.Contains(string(msg), "Failed to install") {
			return m, func() tea.Msg { return setupStepMsg(4) }
		}

	case setupDoneMsg:
		m.setupDone = true
		m.viewState = ViewHome
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		cmds = append(cmds, tickProgress())
		
		if m.notifTimer > 0 {
			m.notifTimer--
			if m.notifTimer <= 0 {
				m.notification = ""
			}
		}
		
		if m.player.State() == player.StatePlaying {
			vizWidth := m.width - 4
			if vizWidth < 1 { vizWidth = 1 }
			m.visualizerBars = make([]int, vizWidth)
			for i := range m.visualizerBars {
				m.visualizerBars[i] = rand.Intn(8)
			}
		} else {
			vizWidth := m.width - 4
			if vizWidth < 1 { vizWidth = 1 }
			m.visualizerBars = make([]int, vizWidth)
		}

		if m.currentTrack != nil && m.player.State() == player.StateStopped && m.statusMsg == "" {
			m.store.AddHistory(*m.currentTrack)
			m.currentTrack = nil
			if track, ok := m.queue.Pop(); ok {
				m.statusMsg = "Resolving next track..."
				cmds = append(cmds, m.resolveCmd(track))
			}
		}
		
		if m.cfg.Radio && m.queue.Len() < m.cfg.QueueSize && m.statusMsg == "" {
			var seedTrack *soundcloud.Track
			if m.currentTrack != nil {
				seedTrack = m.currentTrack
			} else if last, ok := m.queue.LastPlayed(); ok {
				seedTrack = &last
			}
			
			if seedTrack != nil && !m.radioFetched[seedTrack.ID] {
				m.radioFetched[seedTrack.ID] = true
				m.statusMsg = "Fetching related tracks..."
				cmds = append(cmds, m.fetchRelatedCmd(seedTrack.ID))
			}
		}

	case searchResultMsg:
		m.isSearching = false
		m.searchResults = msg
		m.searchCursor = 0
		m.statusMsg = ""
		m.searchInput.Blur()
		cmds = append(cmds, m.checkLyricsBulkCmd(msg))

	case streamResolvedMsg:
		m.statusMsg = ""
		m.currentTrack = &msg.track
		m.currentArtwork = "" // clear previous artwork
		if err := m.player.Play(m.ctx, *msg.stream, msg.track); err != nil {
			m.showNotif("Playback error: " + err.Error())
		}
		cmds = append(cmds, fetchArtworkCmd(msg.track.ID, msg.track.ArtworkURL))
		cmds = append(cmds, m.fetchLyricsCmd(msg.track.Artist, msg.track.Title, msg.track.ID))
		
		// Add to history
		m.store.AddHistory(*m.currentTrack)
		
		// If radio mode is on and we are near the end of queue
		if m.cfg.Radio && m.queue.Len() < 2 {
			if !m.radioFetched[msg.track.ID] {
				m.radioFetched[msg.track.ID] = true
				cmds = append(cmds, m.fetchRelatedCmd(msg.track.ID))
			}
		}

	case lyricsCheckResultMsg:
		m.lyricsChecked[msg.trackID] = true
		if msg.hasLyrics {
			m.lyricsAvailable[msg.trackID] = true
		}

	case artworkMsg:
		if m.currentTrack != nil && m.currentTrack.ID == msg.trackID && msg.art != "" {
			m.currentArtwork = msg.art
		}

	case lyricsMsg:
		if m.currentTrack != nil && m.currentTrack.ID == msg.trackID {
			if msg.lyrics != "" {
				m.currentLyrics = msg.lyrics
				m.hasLyrics = true
			} else {
				m.currentLyrics = ""
				m.hasLyrics = false
			}
		}

	case relatedTracksMsg:
		m.statusMsg = ""
		added := 0
		for _, t := range msg {
			if m.queue.AddIfNotPlayed(t) {
				added++
			}
		}
		if added > 0 {
			m.showNotif("✓ Queue refilled")
		}

	case errMsg:
		m.errorMsg = msg.Error()
		m.statusMsg = ""
		m.isSearching = false
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) searchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.client.SearchTracks(m.ctx, query, 50)
		if err != nil {
			return errMsg(err)
		}
		return searchResultMsg(tracks)
	}
}

func (m *Model) resolveCmd(track soundcloud.Track) tea.Cmd {
	return func() tea.Msg {
		stream, err := m.resolver.Resolve(m.ctx, track, m.cfg.CookieMode)
		if err != nil {
			return errMsg(err)
		}
		return streamResolvedMsg{track: track, stream: stream}
	}
}

func (m *Model) fetchRelatedCmd(trackID int64) tea.Cmd {
	return func() tea.Msg {
		tracks, err := m.client.GetRelatedTracks(m.ctx, trackID, m.cfg.QueueSize)
		if err != nil {
			return errMsg(err)
		}
		return relatedTracksMsg(tracks)
	}
}
