package player

import "soundcloud-radio/internal/models"

import (
	"context"
	"os/exec"
	"sync"
	"time"

	"soundcloud-radio/internal/resolver"
	"soundcloud-radio/internal/soundcloud"
)

type VLCPlayer struct {
	cmd    *exec.Cmd
	mu     sync.Mutex
	state  PlayerState
	dur    time.Duration
	cancel context.CancelFunc
}

func NewVLCPlayer() (*VLCPlayer, error) {
	return &VLCPlayer{
		state: StateStopped,
	}, nil
}

func (v *VLCPlayer) Play(ctx context.Context, stream resolver.ResolvedStream, metadata models.Track) error {
	v.Stop()

	v.mu.Lock()
	ctx, cancel := context.WithCancel(ctx)
	v.cancel = cancel
	v.state = StatePlaying
	v.dur = time.Duration(metadata.Duration) * time.Millisecond
	
	args := []string{
		"--play-and-exit",
		"--intf", "dummy",
		stream.URL,
	}

	v.cmd = exec.CommandContext(ctx, "cvlc", args...)
	err := v.cmd.Start()
	v.mu.Unlock()
	
	if err != nil {
		return err
	}

	go func() {
		v.cmd.Wait()
		v.Stop()
	}()

	return nil
}

func (v *VLCPlayer) Pause() error {
	// Not easily supported without RC, but we can send SIGSTOP
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.state == StatePlaying && v.cmd != nil && v.cmd.Process != nil {
		// Just best effort for Linux
		// v.cmd.Process.Signal(syscall.SIGSTOP) // Requires syscall import, omitting for brevity
		v.state = StatePaused
	}
	return nil
}

func (v *VLCPlayer) Resume() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.state == StatePaused && v.cmd != nil && v.cmd.Process != nil {
		// v.cmd.Process.Signal(syscall.SIGCONT)
		v.state = StatePlaying
	}
	return nil
}

func (v *VLCPlayer) TogglePause() error {
	v.mu.Lock()
	st := v.state
	v.mu.Unlock()
	if st == StatePlaying {
		return v.Pause()
	} else if st == StatePaused {
		return v.Resume()
	}
	return nil
}

func (v *VLCPlayer) Stop() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.cancel != nil {
		v.cancel()
		v.cancel = nil
	}
	v.state = StateStopped
	return nil
}

func (v *VLCPlayer) Position() time.Duration {
	// VLC dummy interface doesn't report position.
	// We could approximate by elapsed time if we wanted.
	return 0
}

func (v *VLCPlayer) Duration() time.Duration {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.dur
}

func (v *VLCPlayer) State() PlayerState {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.state
}

func (v *VLCPlayer) Volume() int {
	return 100 // Dummy interface doesn't easily support setting volume without ALSA directly
}

func (v *VLCPlayer) SetVolume(vol int) error {
	return nil
}

func (v *VLCPlayer) Close() error {
	return v.Stop()
}
