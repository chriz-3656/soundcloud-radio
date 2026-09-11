package ui

import (
	"context"

	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"


	"soundcloud-radio/internal/config"
	"soundcloud-radio/internal/player"
	"soundcloud-radio/internal/queue"
	"soundcloud-radio/internal/resolver"
	"soundcloud-radio/internal/soundcloud"
	"soundcloud-radio/internal/storage"
)

type ViewState int

const (
	ViewSplash ViewState = iota
	ViewSetup
	ViewHome
	ViewSearch
	ViewFavorites
	ViewHistory
	ViewLyrics
	ViewSettings
	ViewHelp
)

type Model struct {
	cfg         *config.Config
	client      *soundcloud.Client
	player      player.Player
	resolver    resolver.Resolver
	queue       *queue.Queue
	store       *storage.Store
	
	viewState   ViewState
	
	width, height int
	
	// Setup state
	setupLogs    []string
	setupDone    bool

	// Search state
	searchInput textinput.Model
	searchResults []soundcloud.Track
	searchCursor  int
	isSearching   bool
	
	// UI State
	queueCursor  int
	favCursor    int
	histCursor   int
	lyricsCursor int
	
	// Notification
	notification string
	notifTimer   int

	// Visualizer
	visualizerBars []int
	
	// Playback state
	currentTrack   *soundcloud.Track
	currentArtwork string
	currentLyrics  string
	hasLyrics      bool
	isPlaying      bool
	
	// Global status
	statusMsg    string
	errorMsg     string
	
	radioFetched map[int64]bool
	lyricsAvailable map[int64]bool
	lyricsChecked   map[int64]bool
	
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewModel(cfg *config.Config, client *soundcloud.Client) *Model {
	p, _ := player.New()
	r := resolver.NewYTDLPResolver()
	q := queue.NewQueue(cfg.HistoryLimit)
	s, _ := storage.NewStore()
	
	ti := textinput.New()
	ti.Placeholder = "Search SoundCloud..."
	ti.Focus()
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Model{
		cfg:          cfg,
		client:       client,
		player:       p,
		resolver:     r,
		queue:        q,
		store:        s,
		viewState:    ViewSplash,
		searchInput:     ti,
		ctx:             ctx,
		cancel:          cancel,
		radioFetched:    make(map[int64]bool),
		lyricsAvailable: make(map[int64]bool),
		lyricsChecked:   make(map[int64]bool),
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
