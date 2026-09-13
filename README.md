# ⚡ Echo Radio TUI

A blazing-fast, premium terminal user interface (TUI) music player. Originally built as a SoundCloud client, Echo Radio has evolved into a dual-provider powerhouse featuring **JioSaavn** (Primary) and **SoundCloud** (Fallback) integration. Built entirely in Go using the Bubble Tea framework.

## 🚀 Features

- **Dual-Provider Ecosystem:** Hot-switch between JioSaavn and SoundCloud using the `Tab` key.
- **Zero-Latency Playback:** JioSaavn streams are natively decrypted (DES-ECB) in Go, completely bypassing `yt-dlp` for instantaneous playback!
- **Dynamic Trending Homepages:** Echo fetches the real-time trending charts for India (JioSaavn) or Top 50 Global (SoundCloud) directly to your home screen.
- **Premium UI:** 3-column layout (Sidebar, Main Content, Now Playing) with a persistent playbar and fully responsive scaling.
- **Vast Search:** Seamlessly search your active provider with up to 50 results at a time, fully paginated.
- **Synchronized Lyrics & Multilingual Support:** Natively integrates the LRCLIB API. Enhanced Unicode `[]rune` parsing guarantees flawless text wrapping for Hindi, Tamil, and other non-ASCII characters without shattering the UI.
- **TrueColor Album Art:** Renders high-fidelity TrueColor (24-bit) album art natively in your terminal with automatic failovers (`500x500` -> `150x150`).
- **Radio Mode (Autoplay):** Let the algorithm take over. Automatically fetches endless related tracks to keep the music playing infinitely based on the current track (fully deduplicated!).
- **Favorites & History:** Persists your listening history and favorite tracks locally.
- **Self-Healing Bootloader:** Automatically detects missing `yt-dlp` dependencies and natively downloads the latest official binary for SoundCloud fallback resolution.

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
go build -o echo-radio cmd/radio/main.go
```

## ⌨️ Keybindings

**Global Commands:**
- `Space` : Play/Pause/Resume
- `n` : Skip to next track
- `p` : Previous track (Pops from history)
- `s` : Stop playback
- `+` / `=` : Volume Up
- `-` / `_` : Volume Down
- `m` : Mute / Unmute
- `Tab` : Switch Provider (JioSaavn <-> SoundCloud)
- `f` / `F` : Favorite / Unfavorite the current track
- `l` / `L` : Toggle Lyrics View
- `/` : Quick search (jumps to search box)
- `q` or `Ctrl+C` : Quit application

**Queue/Home Commands:**
- `a` : Add highlighted track to queue
- `d` : Remove highlighted track from queue
- `c` : Clear the entire queue
- `S` : Shuffle queue

**Navigation:**
- `Up` / `k` : Move cursor up
- `Down` / `j` : Move cursor down
- `Enter` : Select item / Play track
- `Esc` : Blur search box or go back to Home Feed

---

## 📜 Development History & Changelog

### The Journey
This project started as an experiment in building a hyper-responsive Terminal UI for SoundCloud and evolved into a multi-provider feature-rich player named Echo Radio. Over the course of development, we iteratively squashed bugs, integrated complex external APIs, built cryptography layers, and polished the user experience.

### v2.0.0 - v2.1.0: Echo Radio, Dual-Providers, and Zero-Latency ⚡
- **Feature Added (The Rebrand):** Transformed "SoundCloud Radio" into **Echo Radio**. Introduced dynamic ASCII branding and dual-provider support.
- **Feature Added (JioSaavn Native Engine):** Wrote a native Go DES-ECB decryption layer for JioSaavn's `encrypted_media_url`. This bypasses `yt-dlp` completely, dropping resolve latency from 25+ seconds to literal milliseconds!
- **Feature Added (Dynamic Homepages):** Replaced the static queue screen with interactive, provider-specific Home Feeds (e.g., JioSaavn's "Trending in India" charts). Pressing `Tab` instantly swaps the ecosystem, theme, and feeds.
- **Bug Fixed (Non-ASCII Layout Breaks):** Hindi/Tamil song titles were breaking the BubbleTea UI grid. Transitioned all text-slicing logic from standard byte-lengths `len(line)` to true Unicode array slicing `[]rune(line)`.
- **Bug Fixed (Album Art Failovers):** JioSaavn tracks missing high-res (`500x500`) artwork caused standard default ASCII fallback. Engineered automatic fallback degradation to `150x150` native thumbnails.
- **Bug Fixed (Empty Radio Queue):** Auto-queue failed to seed if no track was currently playing. Enter key was hooked to seed the Radio generation pipeline on first playback.
- **Feature Added:** Built API deduplication maps (`seen[id]`) to permanently eliminate duplicate tracks across Search and Radio modules.

### v1.0.1 - v1.0.13: The Foundation & Self-Healing Setup
- **Feature Added:** Massive ASCII startup bootloader that runs system dependency checks. Engineered `EnsureYTDLP()`, an automated self-healing protocol that detects missing binaries and downloads them.
- **Feature Added:** Integrated the **LRCLIB API** for perfectly synchronized lyrics. Implemented Asynchronous Bulk Fetching in the background.
- **Bug Fixed:** `Ctrl+C` was failing due to Bubble Tea's `textinput` interception. Mathematical panics on small terminal screens were resolved by enforcing strict > 85x20 thresholds.

## 👨‍💻 Developer
Developed by **chriz-3656 (Chris Mon Saji)**.
An indie developer & cybersecurity student from India.
Creator of NEURO-RECON, ResuMetric, WebDock, Sky Realms SMP.
