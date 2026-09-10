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
	ViewHome
	ViewSearch
	ViewFavorites
	ViewHistory
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
	
	// Search state
	searchInput textinput.Model
	searchResults []soundcloud.Track
	searchCursor  int
	isSearching   bool
	
	// Queue state
	queueCursor int
	
	// Favorites state
	favCursor int
	
	// History state
	histCursor int

	// Visualizer
	visualizerBars []int
	
	// Playback state
	currentTrack *soundcloud.Track
	isPlaying    bool
	
	// Global status
	statusMsg    string
	errorMsg     string
	notification string
	notifTimer   int
	
	radioFetched map[int64]bool
	
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
		searchInput:  ti,
		ctx:          ctx,
		cancel:       cancel,
		radioFetched: make(map[int64]bool),
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
