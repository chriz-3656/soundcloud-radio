package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"soundcloud-radio/internal/config"
	"soundcloud-radio/internal/soundcloud"
	"soundcloud-radio/ui"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientID, err := soundcloud.FetchClientID(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch client ID: %v\n", err)
		os.Exit(1)
	}

	client := soundcloud.NewClient(cfg.APIBase, clientID)
	
	model := ui.NewModel(cfg, client)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
