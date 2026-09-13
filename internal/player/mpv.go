package player

import "soundcloud-radio/internal/models"

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"soundcloud-radio/internal/resolver"
	
)

type MPVPlayer struct {
	socketPath string
	cmd        *exec.Cmd
	conn       net.Conn
	mu         sync.Mutex
	state      PlayerState
	pos        time.Duration
	dur        time.Duration
	vol        int
	reqID      int
	cancel     context.CancelFunc
}

func NewMPVPlayer() (*MPVPlayer, error) {
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("mpv_ipc_%d.sock", time.Now().UnixNano()))
	return &MPVPlayer{
		socketPath: socketPath,
		state:      StateStopped,
		vol:        100,
	}, nil
}

func (m *MPVPlayer) Play(ctx context.Context, stream resolver.ResolvedStream, metadata models.Track) error {
	m.Stop()

	m.mu.Lock()
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.state = StatePlaying
	m.pos = 0
	m.dur = time.Duration(metadata.Duration) * time.Millisecond
	
	args := []string{
		"--really-quiet",
		"--no-video",
		"--volume=100",
		"--user-agent=Mozilla/5.0",
		fmt.Sprintf("--input-ipc-server=%s", m.socketPath),
		fmt.Sprintf("--title=%s - %s", metadata.Artist, metadata.Title),
		stream.URL,
	}

	m.cmd = exec.CommandContext(ctx, "mpv", args...)
	err := m.cmd.Start()
	m.mu.Unlock()
	
	if err != nil {
		return err
	}

	// Wait for socket to be created
	var conn net.Conn
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err = net.Dial("unix", m.socketPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		m.Stop()
		return fmt.Errorf("failed to connect to mpv IPC: %v", err)
	}

	m.mu.Lock()
	m.conn = conn
	m.mu.Unlock()

	go m.listenIPC(conn)
	
	currentCmd := m.cmd
	go func() {
		currentCmd.Wait()
		m.mu.Lock()
		isActive := (m.cmd == currentCmd)
		m.mu.Unlock()
		if isActive {
			m.Stop()
		}
	}()

	go m.pollPosition(ctx)

	return nil
}

func (m *MPVPlayer) pollPosition(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.sendCommand([]interface{}{"get_property", "time-pos"})
		}
	}
}

func (m *MPVPlayer) listenIPC(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()
		var event struct {
			Event string `json:"event"`
		}
		json.Unmarshal(line, &event)
		
		if event.Event == "playback-restart" {
			// started playing
		}
		
		var response struct {
			Data float64 `json:"data"`
			Err  string  `json:"error"`
		}
		if err := json.Unmarshal(line, &response); err == nil && response.Err == "success" {
			m.mu.Lock()
			m.pos = time.Duration(response.Data * float64(time.Second))
			m.mu.Unlock()
		}
	}
}


func (m *MPVPlayer) sendCommand(cmd []interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn == nil {
		return fmt.Errorf("no active connection")
	}
	
	m.reqID++
	req := map[string]interface{}{
		"command":    cmd,
		"request_id": m.reqID,
	}
	
	b, _ := json.Marshal(req)
	b = append(b, '\n')
	_, err := m.conn.Write(b)
	return err
}

func (m *MPVPlayer) Pause() error {
	m.mu.Lock()
	if m.state == StatePlaying {
		m.state = StatePaused
	}
	m.mu.Unlock()
	return m.sendCommand([]interface{}{"set_property", "pause", true})
}

func (m *MPVPlayer) Resume() error {
	m.mu.Lock()
	if m.state == StatePaused {
		m.state = StatePlaying
	}
	m.mu.Unlock()
	return m.sendCommand([]interface{}{"set_property", "pause", false})
}

func (m *MPVPlayer) TogglePause() error {
	m.mu.Lock()
	if m.state == StatePlaying {
		m.state = StatePaused
	} else if m.state == StatePaused {
		m.state = StatePlaying
	}
	m.mu.Unlock()
	return m.sendCommand([]interface{}{"cycle", "pause"})
}

func (m *MPVPlayer) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	m.state = StateStopped
	os.Remove(m.socketPath)
	return nil
}

func (m *MPVPlayer) Position() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Ideally we query IPC periodically, but for simplicity let's just return what we have
	// or we can query it on demand. Since querying is async, we'll start a goroutine in Play to poll position
	return m.pos
}

func (m *MPVPlayer) Duration() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dur
}

func (m *MPVPlayer) State() PlayerState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

func (m *MPVPlayer) Volume() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.vol
}

func (m *MPVPlayer) SetVolume(vol int) error {
	m.mu.Lock()
	if vol < 0 { vol = 0 }
	if vol > 130 { vol = 130 }
	m.vol = vol
	m.mu.Unlock()
	return m.sendCommand([]interface{}{"set_property", "volume", vol})
}

func (m *MPVPlayer) Close() error {
	return m.Stop()
}
