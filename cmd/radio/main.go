package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"soundcloud-radio/internal/config"
	"soundcloud-radio/internal/provider"
	"soundcloud-radio/internal/jiosaavn"
	"soundcloud-radio/internal/soundcloud"
	"soundcloud-radio/ui"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize JioSaavn (Primary)
	jioClient := jiosaavn.NewClient()

	// Initialize SoundCloud (Fallback)
	scClientID, err := soundcloud.FetchClientID(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch soundcloud client ID: %v\n", err)
	}
	scClient := soundcloud.NewClient(cfg.APIBase, scClientID)
	
	providers := []provider.Provider{jioClient, scClient}
	
	model := ui.NewModel(cfg, providers)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
