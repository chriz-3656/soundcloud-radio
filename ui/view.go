package ui

import (
	"fmt"
	"strings"
	"time"

	tea 	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	
	"soundcloud-radio/internal/models"
	"soundcloud-radio/internal/player"
)

var (
	brandStyle        lipgloss.Style
	titleStyle        lipgloss.Style
	artistStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	infoStyle         lipgloss.Style
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	statusStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("77"))
	selectedItemStyle lipgloss.Style
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	mutedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	
	panelBorderStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238"))
	activeBorderStyle lipgloss.Style
)

func UpdateTheme(providerName string) {
	if providerName == "JioSaavn" {
		brandStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("43")).Bold(true) // Cyan-Green
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
		infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
		selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("43")).Bold(true).Background(lipgloss.Color("236"))
		activeBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("43"))
	} else { // SoundCloud
		brandStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Orange
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
		infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
		selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true).Background(lipgloss.Color("236"))
		activeBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("208"))
	}
}

func init() {
	UpdateTheme("JioSaavn") // Default theme
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func formatArtistAlbum(t models.Track) string {
	if t.Album != "" {
		return t.Artist + " • " + t.Album
	}
	return t.Artist
}

func renderArtwork() string {
	return mutedStyle.Render("┌─────────────┐\n│      ♪      │\n│ SOUNDCLOUD  │\n│             │\n└─────────────┘")
}

func (m Model) View() tea.View {
	if m.viewState == ViewSplash {
		splashArt := `
  ██████  ██████ ██       ██████  ██    ██ ██████  
 ██      ██      ██      ██    ██ ██    ██ ██   ██ 
 ███████ ██      ██      ██    ██ ██    ██ ██   ██ 
      ██ ██      ██      ██    ██ ██    ██ ██   ██ 
 ██████  ██████  ███████  ██████   ██████  ██████  
`
		content := brandStyle.Render(splashArt) + "\n\n" + itemStyle.Render("Press Enter to open the main UI")
		v := tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content))
		v.AltScreen = true
		return v
	}

	if m.viewState == ViewSetup {
		var sb strings.Builder
		sb.WriteString(titleStyle.Render("SYSTEM INITIALIZATION") + "\n\n")
		for _, log := range m.setupLogs {
			if strings.Contains(log, "[✓]") {
				sb.WriteString(successStyle.Render(log) + "\n")
			} else if strings.Contains(log, "[!]") || strings.Contains(log, "[✕]") {
				sb.WriteString(errorStyle.Render(log) + "\n")
			} else {
				sb.WriteString(itemStyle.Render(log) + "\n")
			}
		}
		
		v := tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, panelBorderStyle.Width(60).Height(20).Render(sb.String())))
		v.AltScreen = true
		return v
	}

	if m.width < 85 || m.height < 20 {
		// Small terminal fallback
		var sb strings.Builder
		sb.WriteString(brandStyle.Render("SOUNDCLOUD RADIO") + "\n\n")
		if m.currentTrack != nil {
			sb.WriteString(titleStyle.Render("▶ " + m.currentTrack.Title) + "\n")
			sb.WriteString(artistStyle.Render("  " + m.currentTrack.Artist) + "\n\n")
			
			pos := m.player.Position()
			dur := m.player.Duration()
			sb.WriteString(fmt.Sprintf("%s / %s\n\n", formatDuration(pos), formatDuration(dur)))
		} else {
			sb.WriteString(mutedStyle.Render("No track playing.\n\n"))
		}
		sb.WriteString(mutedStyle.Render("[Space] Play/Pause   [/] Search   [q] Quit\n"))
		
		v := tea.NewView(sb.String())
		v.AltScreen = true
		return v
	}

	sidebarWidth := 22
	playbarHeight := 4
	nowPlayingWidth := 35
	mainWidth := m.width - sidebarWidth - nowPlayingWidth - 6
	mainHeight := m.height - playbarHeight - 2

	// === 1. SIDEBAR ===
	var sb strings.Builder
	sb.WriteString(brandStyle.Render("♫ SOUNDCLOUD") + "\n\n")
	
	menuItems := []struct{
		label string
		state ViewState
	}{
		{"  HOME", ViewHome},
		{"  QUEUE", ViewQueue},
		{"  SEARCH", ViewSearch},
		{"  FAVORITES", ViewFavorites},
		{"  HISTORY", ViewHistory},
		{"  ", ViewState(-1)},
		{"  HELP", ViewHelp},
		{"  SETTINGS", ViewSettings},
	}
	
	for _, item := range menuItems {
		if item.state == ViewState(-1) {
			sb.WriteString("\n")
			continue
		}
		if m.viewState == item.state {
			sb.WriteString(selectedItemStyle.Width(sidebarWidth-2).Render("▶ " + item.label) + "\n\n")
		} else {
			sb.WriteString(itemStyle.Width(sidebarWidth-2).Render("  " + item.label) + "\n\n")
		}
	}
	sb.WriteString("\n")
	
	sidebar := panelBorderStyle.Width(sidebarWidth).Height(mainHeight).Render(sb.String())

	// === 2. MAIN CONTENT ===
	var mb strings.Builder
	
	if m.errorMsg != "" {
		mb.WriteString(errorStyle.Render("✕ Error: " + m.errorMsg) + "\n\n")
	} else if m.notification != "" {
		mb.WriteString(successStyle.Render(m.notification) + "\n\n")
	} else if m.statusMsg != "" {
		mb.WriteString(statusStyle.Render("◌ " + m.statusMsg) + "\n\n")
	} else {
		mb.WriteString("\n\n") // Spacing
	}

	switch m.viewState {
	case ViewHome:
		mb.WriteString(titleStyle.Render(strings.ToUpper(m.activeProvider.GetName()) + " HOME") + "\n\n")
		
		if len(m.homeFeed) == 0 {
			if m.errorMsg != "" {
				mb.WriteString(errorStyle.Render(m.errorMsg) + "\n")
			} else {
				mb.WriteString(mutedStyle.Render("Loading trending feed...\n"))
			}
		} else {
			itemsPerPage := (mainHeight - 10) / 2
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.homeCursor >= itemsPerPage {
				startIndex = m.homeCursor - itemsPerPage + 1
			}
			endIndex := startIndex + itemsPerPage
			if endIndex > len(m.homeFeed) { endIndex = len(m.homeFeed) }

			for i := startIndex; i < endIndex; i++ {
				t := m.homeFeed[i]
				cursor := "  "
				style := itemStyle
				if i == m.homeCursor {
					cursor = "▶ "
					style = selectedItemStyle
				}
				
				lyricsIcon := ""
				if m.lyricsAvailable[t.ID] {
					lyricsIcon = " 📜"
				}
				
				line := fmt.Sprintf("%s%02d  %s%s", cursor, i+1, t.Title, lyricsIcon)
				runes := []rune(line)
				if len(runes) > mainWidth - 4 {
					line = string(runes[:mainWidth-7]) + "..."
				}
				mb.WriteString(style.Width(mainWidth-2).Render(line) + "\n")
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("      %s", formatArtistAlbum(t))) + "\n")
			}
			if len(m.homeFeed) > itemsPerPage {
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("\n  ... %d trending tracks", len(m.homeFeed))))
			}
		}
		
	case ViewQueue:
		mb.WriteString(titleStyle.Render("UP NEXT") + "\n\n")
		items := m.queue.Items()
		if len(items) == 0 {
			mb.WriteString(mutedStyle.Render("Queue is empty."))
		} else {
			itemsPerPage := mainHeight - 8
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.queueCursor >= itemsPerPage {
				startIndex = m.queueCursor - itemsPerPage + 1
			}
			endIndex := startIndex + itemsPerPage
			if endIndex > len(items) { endIndex = len(items) }
			
			for i := startIndex; i < endIndex; i++ {
				t := items[i]
				cursor := "  "
				style := itemStyle
				if i == m.queueCursor {
					cursor = "▶ "
					style = selectedItemStyle
				}
				
				lyricsIcon := ""
				if m.lyricsAvailable[t.ID] {
					lyricsIcon = " 📜"
				}
				
				line := fmt.Sprintf("%s%02d  %s%s", cursor, i+1, t.Title, lyricsIcon)
				runes := []rune(line)
				if len(runes) > mainWidth - 4 {
					line = string(runes[:mainWidth-7]) + "..."
				}
				mb.WriteString(style.Width(mainWidth-2).Render(line) + "\n")
			}
			if len(items) > itemsPerPage {
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("  ... %d tracks total", len(items))))
			}
		}
		
	case ViewSearch:
		mb.WriteString(titleStyle.Render("SEARCH SOUNDCLOUD") + "\n\n")
		if m.searchInput.Focused() {
			mb.WriteString("Query: " + m.searchInput.View() + "\n\n")
		} else {
			mb.WriteString("Query: " + titleStyle.Render(m.searchInput.Value()) + "\n\n")
		}
		
		if m.isSearching {
			mb.WriteString(statusStyle.Render("Searching...\n"))
		} else {
			itemsPerPage := (mainHeight - 10) / 2
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.searchCursor >= itemsPerPage {
				startIndex = m.searchCursor - itemsPerPage + 1
			}
			endIndex := startIndex + itemsPerPage
			if endIndex > len(m.searchResults) { endIndex = len(m.searchResults) }

			for i := startIndex; i < endIndex; i++ {
				t := m.searchResults[i]
				cursor := "  "
				style := itemStyle
				if i == m.searchCursor {
					cursor = "▶ "
					style = selectedItemStyle
				}
				
				lyricsIcon := ""
				if m.lyricsAvailable[t.ID] {
					lyricsIcon = " 📜"
				}
				
				line := fmt.Sprintf("%s%02d  %s%s", cursor, i+1, t.Title, lyricsIcon)
				runes := []rune(line)
				if len(runes) > mainWidth - 4 {
					line = string(runes[:mainWidth-7]) + "..."
				}
				mb.WriteString(style.Width(mainWidth-2).Render(line) + "\n")
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("      %s", formatArtistAlbum(t))) + "\n")
			}
			if len(m.searchResults) > itemsPerPage {
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("\n  ... %d results total", len(m.searchResults))))
			}
		}
		
	case ViewFavorites:
		mb.WriteString(titleStyle.Render("FAVORITES") + "\n\n")
		favs := m.store.GetFavorites()
		if len(favs) == 0 {
			mb.WriteString(mutedStyle.Render("No favorites yet. Press 'f' to add tracks!"))
		} else {
			itemsPerPage := mainHeight - 8
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.favCursor >= itemsPerPage { startIndex = m.favCursor - itemsPerPage + 1 }
			endIndex := startIndex + itemsPerPage
			if endIndex > len(favs) { endIndex = len(favs) }

			for i := startIndex; i < endIndex; i++ {
				t := favs[i]
				cursor := "  "
				style := itemStyle
				if i == m.favCursor {
					cursor = "▶ "
					style = selectedItemStyle
				}
				
				lyricsIcon := ""
				if m.lyricsAvailable[t.ID] {
					lyricsIcon = " 📜"
				}
				
				line := fmt.Sprintf("%s%s%s", cursor, t.Title, lyricsIcon)
				runes := []rune(line)
				if len(runes) > mainWidth - 4 {
					line = string(runes[:mainWidth-7]) + "..."
				}
				mb.WriteString(style.Width(mainWidth-2).Render(line) + "\n")
			}
		}
		
	case ViewHistory:
		mb.WriteString(titleStyle.Render("HISTORY") + "\n\n")
		hist := m.store.GetHistory()
		if len(hist) == 0 {
			mb.WriteString(mutedStyle.Render("No history."))
		} else {
			itemsPerPage := mainHeight - 8
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.histCursor >= itemsPerPage { startIndex = m.histCursor - itemsPerPage + 1 }
			endIndex := startIndex + itemsPerPage
			if endIndex > len(hist) { endIndex = len(hist) }

			for i := startIndex; i < endIndex; i++ {
				h := hist[i]
				cursor := "  "
				style := itemStyle
				if i == m.histCursor {
					cursor = "▶ "
					style = selectedItemStyle
				}
				
				lyricsIcon := ""
				if m.lyricsAvailable[h.Track.ID] {
					lyricsIcon = " 📜"
				}
				
				line := fmt.Sprintf("%s%s%s", cursor, h.Track.Title, lyricsIcon)
				runes := []rune(line)
				if len(runes) > mainWidth - 4 {
					line = string(runes[:mainWidth-7]) + "..."
				}
				mb.WriteString(style.Width(mainWidth-2).Render(line) + "\n")
			}
		}

	case ViewLyrics:
		mb.WriteString(titleStyle.Render("LYRICS") + "\n\n")
		if m.currentLyrics == "" {
			mb.WriteString(mutedStyle.Render("No lyrics available."))
		} else {
			lines := strings.Split(m.currentLyrics, "\n")
			itemsPerPage := mainHeight - 8
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.lyricsCursor >= itemsPerPage { startIndex = m.lyricsCursor - itemsPerPage + 1 }
			endIndex := startIndex + itemsPerPage
			if endIndex > len(lines) { endIndex = len(lines) }

			for i := startIndex; i < endIndex; i++ {
				l := lines[i]
				style := itemStyle
				if i == m.lyricsCursor {
					style = selectedItemStyle
				}
				mb.WriteString(style.Width(mainWidth-2).Render(l) + "\n")
			}
		}
		
	case ViewHelp:
		mb.WriteString(titleStyle.Render("HELP & COMMANDS") + "\n\n")
		helpTxt := `PLAYBACK
  Space    Play/Pause
  n        Next
  s        Stop
  + / -    Volume Up/Down

NAVIGATION
  ↑ / k    Up
  ↓ / j    Down
  Enter    Select
  Esc      Back
  /        Search

QUEUE
  a        Add
  d        Remove
  c        Clear
  S        Shuffle
  
GLOBAL
  r        Toggle Radio Mode
  f        Favorite track
  q        Quit`
		mb.WriteString(itemStyle.Render(helpTxt))
		
	case ViewSettings:
		mb.WriteString(titleStyle.Render("SETTINGS & ABOUT") + "\n\n")
		mb.WriteString(itemStyle.Render("Player: mpv/vlc auto-detected\n"))
		mb.WriteString(itemStyle.Render(fmt.Sprintf("Volume: %d%%\n", m.player.Volume())))
		mb.WriteString(itemStyle.Render(fmt.Sprintf("Radio: %v\n", m.cfg.Radio)))
		mb.WriteString(itemStyle.Render(fmt.Sprintf("Cookie Mode: %s\n\n", m.cfg.CookieMode)))
		
		mb.WriteString(titleStyle.Render("DEVELOPER INFO") + "\n")
		devInfo := `Developed by chriz-3656 (Chris Mon Saji)
An indie developer & cybersecurity student from India.
Projects: NEURO-RECON, ResuMetric, WebDock, Sky Realms SMP.`
		mb.WriteString(mutedStyle.Render(devInfo))
	}
	
	mainPanelStyle := panelBorderStyle
	if m.searchInput.Focused() {
		mainPanelStyle = activeBorderStyle
	}
	mainContent := mainPanelStyle.Width(mainWidth).Height(mainHeight).Render(mb.String())

	// === 3. NOW PLAYING PANEL ===
	var np strings.Builder
	np.WriteString(titleStyle.Render("NOW PLAYING") + "\n\n")
	
	if m.currentArtwork != "" {
		np.WriteString(m.currentArtwork + "\n\n")
	} else {
		np.WriteString(renderArtwork() + "\n\n")
	}
	
	if m.currentTrack != nil {
		t := m.currentTrack.Title
		if len(t) > nowPlayingWidth-4 { t = t[:nowPlayingWidth-7] + "..." }
		np.WriteString(titleStyle.Render(t) + "\n")
		
		a := m.currentTrack.Artist
		if len(a) > nowPlayingWidth-4 { a = a[:nowPlayingWidth-7] + "..." }
		np.WriteString(artistStyle.Render(a) + "\n\n")
		
		pos := m.player.Position()
		dur := m.player.Duration()
		
		// Progress bar
		barWidth := nowPlayingWidth - 6
		var percent float64
		if dur > 0 {
			percent = float64(pos) / float64(dur)
		}
		filled := int(float64(barWidth) * percent)
		if filled < 0 { filled = 0 }
		if filled > barWidth { filled = barWidth }
		empty := barWidth - filled
		bar := strings.Repeat("━", filled) + "●" + strings.Repeat("━", empty)
		
		np.WriteString(infoStyle.Render(bar) + "\n")
		np.WriteString(mutedStyle.Render(fmt.Sprintf("%s / %s", formatDuration(pos), formatDuration(dur))) + "\n\n")
		
		if m.store.IsFavorite(m.currentTrack.ID) {
			np.WriteString(brandStyle.Render("♥ Favorited") + "\n\n")
		} else {
			np.WriteString("\n")
		}
		
		if m.hasLyrics {
			np.WriteString(infoStyle.Render("[l] Show Lyrics") + "\n\n")
		} else {
			np.WriteString("\n\n")
		}
		
	} else {
		np.WriteString(mutedStyle.Render("Nothing is playing."))
	}

	// Visualizer block
	vizChars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	var vizBuilder strings.Builder
	for _, val := range m.visualizerBars {
		if val < 0 { val = 0 }
		if val > 7 { val = 7 }
		vizBuilder.WriteString(vizChars[val])
	}
	vizStr := statusStyle.Render(vizBuilder.String())

	// Push visualizer to the bottom of the panel
	npTop := lipgloss.NewStyle().Height(mainHeight - 2).Render(np.String())
	npFull := lipgloss.JoinVertical(lipgloss.Top, npTop, vizStr)

	nowPlayingContent := panelBorderStyle.Width(nowPlayingWidth).Height(mainHeight).Render(npFull)

	// === 4. PLAYBAR (BOTTOM) ===
	var pb strings.Builder
	
	stateIcon := "▶"
	if m.player.State() == player.StatePaused {
		stateIcon = "Ⅱ"
	} else if m.player.State() == player.StateStopped {
		stateIcon = "■"
	}
	
	radioStatus := mutedStyle.Render("RADIO OFF")
	if m.cfg.Radio {
		radioStatus = brandStyle.Render("RADIO ON")
	}
	
	volStr := fmt.Sprintf("🔊 %d%%", m.player.Volume())
	queueStr := fmt.Sprintf("QUEUE %d", m.queue.Len())
	
	pb.WriteString(fmt.Sprintf("%s %s    %s    %s    %s    STREAM ●\n", 
		statusStyle.Render(stateIcon), 
		statusStyle.Render(string(m.player.State())), 
		itemStyle.Render(volStr),
		radioStatus,
		itemStyle.Render(queueStr),
	))
	
	// Dynamic Footer Shortcuts
	shortcuts := "[Space] Play/Pause   [n] Next   [s] Stop   [/] Search   [?] Help   [q] Quit"
	if m.viewState == ViewSearch {
		shortcuts = "[Enter] Play/Select   [a] Add to Queue   [r] Radio Mode   [Esc] Back"
	} else if m.viewState == ViewHome { // Queue
		shortcuts = "[Enter] Play   [a] Add   [d] Remove   [c] Clear   [S] Shuffle   [/] Search"
	}
	pb.WriteString(mutedStyle.Render(shortcuts))
	
	playbar := panelBorderStyle.Width(m.width - 2).Height(playbarHeight - 2).Render(pb.String())

	// === ASSEMBLE ===
	topSection := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainContent, nowPlayingContent)
	fullUI := lipgloss.JoinVertical(lipgloss.Left, topSection, playbar)
	
	v := tea.NewView(fullUI)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
