package ui

import "soundcloud-radio/internal/models"

import (
	"context"

	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"


	"soundcloud-radio/internal/config"
	"soundcloud-radio/internal/player"
	"soundcloud-radio/internal/queue"
	"soundcloud-radio/internal/resolver"
	"soundcloud-radio/internal/provider"
	
	"soundcloud-radio/internal/storage"
)

type ViewState int

const (
	ViewSplash ViewState = iota
	ViewSetup
	ViewHome
	ViewQueue
	ViewSearch
	ViewFavorites
	ViewHistory
	ViewLyrics
	ViewSettings
	ViewHelp
)

type Model struct {
	cfg            *config.Config
	providers      []provider.Provider
	providerIndex  int
	activeProvider provider.Provider
	player         player.Player
	resolver       resolver.Resolver
	queue          *queue.Queue
	store          *storage.Store
	
	viewState   ViewState
	
	width, height int
	
	// Setup state
	setupLogs    []string
	setupDone    bool

	// Search state
	searchInput textinput.Model
	searchResults []models.Track
	searchCursor  int
	isSearching   bool
	
	// UI State
	queueCursor  int
	favCursor    int
	histCursor   int
	lyricsCursor int
	homeCursor   int
	homeFeed     []models.Track
	
	// Notification
	notification string
	notifTimer   int

	// Visualizer
	visualizerBars []int
	
	// Playback state
	currentTrack   *models.Track
	currentArtwork string
	currentLyrics  string
	hasLyrics      bool
	isPlaying      bool
	
	// Global status
	statusMsg    string
	errorMsg     string
	
	radioFetched    map[string]bool
	lyricsAvailable map[string]bool
	lyricsChecked   map[string]bool
	
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewModel(cfg *config.Config, providers []provider.Provider) *Model {
	p, _ := player.New()
	r := resolver.NewYTDLPResolver()
	q := queue.NewQueue(cfg.HistoryLimit)
	s, _ := storage.NewStore()
	
	ti := textinput.New()
	ti.Placeholder = "Search " + providers[0].GetName() + "..."
	ti.Focus()
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Model{
		cfg:             cfg,
		providers:       providers,
		providerIndex:   0,
		activeProvider:  providers[0],
		player:          p,
		resolver:        r,
		queue:           q,
		store:           s,
		viewState:       ViewSplash,
		searchInput:     ti,
		ctx:             ctx,
		cancel:          cancel,
		radioFetched:    make(map[string]bool),
		lyricsAvailable: make(map[string]bool),
		lyricsChecked:   make(map[string]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tickProgress(),
	)
}

type tickMsg time.Time
func tickProgress() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type homeFeedMsg []models.Track

func (m *Model) fetchHomeFeedCmd() tea.Cmd {
	return func() tea.Msg {
		feed, err := m.activeProvider.GetHomeFeed(m.ctx)
		if err != nil {
			return errMsg(err)
		}
		return homeFeedMsg(feed)
	}
}
