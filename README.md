# SoundCloud Radio TUI

A blazing-fast, premium terminal user interface (TUI) music player for SoundCloud, designed with a modern architecture inspired by Spotify-TUI and YouTube Music. Built entirely in Go using the Bubble Tea framework.

## 🚀 Features

- **Premium UI:** 3-column layout (Sidebar, Main Content, Now Playing) with a persistent playbar.
- **Vast Search:** Seamlessly search SoundCloud with up to 50 results at a time, fully paginated.
- **Synchronized Lyrics:** Natively integrates the LRCLIB API to fetch and display synchronized lyrics for tracks (with a fast `📜` availability indicator).
- **TrueColor Album Art:** Renders high-fidelity TrueColor (24-bit) album art using ANSI blocks directly in your terminal.
- **Radio Mode (Autoplay):** Let the algorithm take over. Automatically fetches related tracks to keep the music playing infinitely based on your queue history.
- **Favorites & History:** Persists your listening history and favorite tracks locally to `~/.config/soundcloud-radio/data.json`.
- **Intelligent Resizing:** Gracefully falls back to a compact, text-only layout on small terminals (< 85x20).
- **Audio Engine Integration:** Direct IPC integration with `mpv` for zero-latency playback control, volume adjustments, and exact progress tracking.
- **Self-Healing Bootloader:** Automatically detects missing `yt-dlp` dependencies and natively downloads the latest official binary to the application cache during the boot sequence.
- **Mouse Support:** Full mouse clicking support for seamless sidebar navigation and search focus.

## 📦 Installation

Download the pre-compiled binary from the [Releases](https://github.com/chriz-3656/soundcloud-radio/releases) page for your architecture (Linux AMD64 / ARM64).

### Dependencies

Ensure you have the following installed on your system:
- `mpv` (Required for audio backend)

*Note: The application will automatically download `yt-dlp` in the background if it is not installed globally.*

### Building from Source

```bash
git clone https://github.com/chriz-3656/soundcloud-radio.git
cd soundcloud-radio
make build
```
The binary will be located at `bin/soundcloud-radio`.

## ⌨️ Keybindings

**Global Commands:**
- `Space` : Play/Pause/Resume
- `n` : Skip to next track
- `p` : Previous track (Pops from history)
- `s` : Stop playback
- `+` / `=` : Volume Up
- `-` / `_` : Volume Down
- `m` : Mute / Unmute
- `f` / `F` : Favorite / Unfavorite the current track
- `l` / `L` : Toggle Lyrics View
- `/` : Quick search (jumps to search box)
- `q` or `Ctrl+C` : Quit application

**Queue Commands:**
- `a` : Add highlighted track to queue
- `d` : Remove highlighted track from queue
- `c` : Clear the entire queue
- `S` : Shuffle queue

**Navigation:**
- `Up` / `k` : Move cursor up
- `Down` / `j` : Move cursor down
- `Enter` : Select item / Play track
- `Esc` : Blur search box or go back to Queue

---

## 📜 Development History & Changelog

### The Journey
This project started as an experiment in building a hyper-responsive Terminal UI for SoundCloud and evolved into a feature-rich, standalone music player. Over the course of development, we iteratively squashed bugs, integrated complex external APIs, and polished the user experience.

### v1.0.1 - v1.0.6: The Foundation & Self-Healing Setup
- **Feature Added:** Implemented a massive ASCII `SCLOUD.` startup bootloader that runs system dependency checks.
- **Bug Fixed:** "All cookie strategies failed to resolve stream".
  - *Cause:* Swallowing standard errors in `yt-dlp` execution masked the fact that `yt-dlp` wasn't actually installed on the system.
  - *Solution:* Engineered `EnsureYTDLP()`, an automated self-healing protocol that detects if `yt-dlp` is missing, downloads the latest binary directly from GitHub, and executes it from `~/.config/soundcloud-radio/`.
- **Bug Fixed:** `Ctrl+C` was not quitting the app when searching.
  - *Cause:* Bubble Tea's `textinput` component was intercepting global OS kill sequences.
  - *Solution:* Added an exclusion clause to bypass standard kill sequences when the input field is focused.
- **Bug Fixed:** Album Art reverting to generic SoundCloud logos.
  - *Cause:* Attempting to fetch `-large.jpg` for album art would result in a `404 Not Found` for certain tracks, causing the engine to fall back to a default logo.
  - *Solution:* Rewrote the resolver to pull the native `-t200x200.jpg` CDN URLs before scaling them down into terminal ANSI grids.

### v1.0.7 - v1.0.9: Lyrics Integration & Keybinding Polish
- **Feature Added:** Integrated the **LRCLIB API** for perfectly synchronized lyrics. Pressing `l` or `L` dynamically replaces the main interface with a scrollable lyrics sheet.
- **Bug Fixed:** LRCLIB API returning `503 Service Unavailable` or `context deadline exceeded`.
  - *Cause:* Network latency causing the default 5s HTTP client timeout to trigger early.
  - *Solution:* Bumped HTTP timeout to 10s. Identified 30-second SoundCloud Go+ paywall restrictions for premium tracks.
- **Bug Fixed:** Uppercase `L` and `F` inputs were being ignored by the key listener.
  - *Cause:* Exact character-string matching in Bubble Tea event loop.
  - *Solution:* Hardcoded explicit string routing for upper-case keystrokes.

### v1.0.10 - v1.0.13: UI Refinement & Async Bulk Fetching
- **Bug Fixed:** Mathematical `panic()` when terminal was resized below 63 columns.
  - *Cause:* Subtracting fixed UI layout integers from a terminal width that was too small resulted in a negative width constraint being passed to Lipgloss layout renderers.
  - *Solution:* Raised the safety fallback threshold to `85x20` to guarantee the engine collapses to the small-terminal ASCII UI gracefully.
- **Feature Added:** Built an **Asynchronous Bulk LRCLIB Fetcher**. When searching, the engine silently fires off 50 concurrent requests in a non-blocking background pool. If lyrics are discovered, a `📜` emoji seamlessly pops into the UI adjacent to the track.
- **Feature Added:** Extracted the Audio Visualizer from the Now Playing sidebar and spanned it dynamically across 31 blocks specifically locked to the bottom of the right-hand panel for extreme visual flair.

## 👨‍💻 Developer
Developed by **chriz-3656 (Chris Mon Saji)**.
An indie developer & cybersecurity student from India.
Creator of NEURO-RECON, ResuMetric, WebDock, Sky Realms SMP.
