package player

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"soundcloud-radio/internal/resolver"
	"soundcloud-radio/internal/soundcloud"
)

type PlayerState string

const (
	StateStopped PlayerState = "Stopped"
	StatePlaying PlayerState = "Playing"
	StatePaused  PlayerState = "Paused"
)

type Player interface {
	Play(ctx context.Context, stream resolver.ResolvedStream, metadata soundcloud.Track) error
	Pause() error
	Resume() error
	TogglePause() error
	Stop() error
	Position() time.Duration
	Duration() time.Duration
	State() PlayerState
	Volume() int
	SetVolume(vol int) error
	Close() error
}

func New() (Player, error) {
	if _, err := os.LookupEnv("FORCE_VLC"); err == false {
		if p, err := exec.LookPath("mpv"); err == nil && p != "" {
			return NewMPVPlayer()
		}
	}
	if p, err := exec.LookPath("cvlc"); err == nil && p != "" {
		return NewVLCPlayer()
	}
	if p, err := exec.LookPath("vlc"); err == nil && p != "" {
		return NewVLCPlayer()
	}
	return nil, fmt.Errorf("no media player found (mpv, cvlc, vlc)")
}

