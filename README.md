# SoundCloud Radio TUI

A blazing-fast, premium terminal user interface (TUI) music player for SoundCloud, designed with a modern architecture inspired by Spotify-TUI and YouTube Music. Built entirely in Go using the Bubble Tea framework.

## 🚀 Features

- **Premium UI:** 3-column layout (Sidebar, Main Content, Now Playing) with a persistent playbar.
- **Vast Search:** Seamlessly search SoundCloud with up to 50 results at a time, fully paginated.
- **Radio Mode (Autoplay):** Let the algorithm take over. Automatically fetches related tracks to keep the music playing infinitely based on your queue history.
- **Favorites & History:** Persists your listening history and favorite tracks locally to `~/.config/soundcloud-radio/data.json`.
- **Intelligent Resizing:** Gracefully falls back to a compact, text-only layout on small terminals (< 60x15).
- **Audio Engine Integration:** Direct IPC integration with `mpv` for zero-latency playback control, volume adjustments, and exact progress tracking. Automatically falls back to `vlc` if `mpv` isn't installed.
- **Mouse Support:** Full mouse clicking support for seamless sidebar navigation and search focus.

## 📦 Installation

Download the pre-compiled binary from the [Releases](https://github.com/chriz-3656/soundcloud-radio/releases) page for your architecture (Linux AMD64 / ARM64).

### Dependencies

Ensure you have the following installed on your system:
- `yt-dlp` (Required for resolving stream URLs)
- `mpv` (Highly recommended for optimal playback & IPC controls)
- `vlc` (Optional fallback if mpv is unavailable)

### Building from Source

```bash
git clone https://github.com/chriz-3656/soundcloud-radio.git
cd soundcloud-radio
make build
```
The binary will be located at `bin/soundcloud-radio`.

## ⌨️ Keybindings

**Global Commands (Available everywhere):**
- `Space` : Play/Pause/Resume
- `n` : Skip to next track
- `p` : Reback to previous track (Pops from history)
- `s` : Stop playback
- `+` / `=` : Volume Up
- `-` / `_` : Volume Down
- `m` : Mute / Unmute
- `f` : Favorite / Unfavorite the current track
- `/` : Quick search (jumps to search box)
- `q` or `Ctrl+C` : Quit application

**Queue Commands:**
- `a` : Add highlighted track to queue
- `d` : Remove highlighted track from queue
- `c` : Clear the entire queue
- `S` (Shift+S) : Shuffle queue

**Navigation:**
- `Up` / `k` : Move cursor up
- `Down` / `j` : Move cursor down
- `Enter` : Select item / Play track
- `Esc` : Blur search box or go back to Queue

## 🛠️ Configuration
The player automatically stores its data in `~/.config/soundcloud-radio/`.

## 👨‍💻 Developer
Developed by **chriz-3656 (Chris Mon Saji)**.
An indie developer & cybersecurity student from India.
Creator of NEURO-RECON, ResuMetric, WebDock, Sky Realms SMP.
