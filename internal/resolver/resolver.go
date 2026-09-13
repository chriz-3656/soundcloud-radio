package resolver

import "soundcloud-radio/internal/models"

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	
)

type ResolvedStream struct {
	URL string
}

type Resolver interface {
	Resolve(ctx context.Context, track models.Track, cookieMode string) (*ResolvedStream, error)
}

type YTDLPResolver struct{}

func NewYTDLPResolver() *YTDLPResolver {
	return &YTDLPResolver{}
}

func (r *YTDLPResolver) Resolve(ctx context.Context, track models.Track, cookieMode string) (*ResolvedStream, error) {
	if cookieMode == "auto" {
		var errs []string
		
		// Try firefox
		url, err := r.resolveWithCookies(ctx, track.PermalinkURL, "firefox")
		if err == nil && url != "" { return &ResolvedStream{URL: url}, nil }
		if err != nil { errs = append(errs, "firefox: "+err.Error()) }

		// Try chrome
		url, err = r.resolveWithCookies(ctx, track.PermalinkURL, "chrome")
		if err == nil && url != "" { return &ResolvedStream{URL: url}, nil }
		if err != nil { errs = append(errs, "chrome: "+err.Error()) }

		// Try none
		url, err = r.resolveWithCookies(ctx, track.PermalinkURL, "")
		if err == nil && url != "" { return &ResolvedStream{URL: url}, nil }
		if err != nil { errs = append(errs, "none: "+err.Error()) }

		return nil, fmt.Errorf("resolution failed: %s", strings.Join(errs, " | "))
	} else if cookieMode != "none" && cookieMode != "" {
		url, err := r.resolveWithCookies(ctx, track.PermalinkURL, cookieMode)
		if err == nil && url != "" {
			return &ResolvedStream{URL: url}, nil
		}
		return nil, err
	}

	url, err := r.resolveWithCookies(ctx, track.PermalinkURL, "")
	if err == nil && url != "" {
		return &ResolvedStream{URL: url}, nil
	}
	return nil, err
}

func (r *YTDLPResolver) resolveWithCookies(ctx context.Context, trackURL, cookies string) (string, error) {
	binPath, err := EnsureYTDLP(ctx)
	if err != nil {
		return "", err
	}

	args := []string{"-g", "-f", "bestaudio", "--no-warnings", "--no-playlist"}
	if cookies != "" {
		args = append(args, "--cookies-from-browser", cookies)
	}
	args = append(args, trackURL)

	cmd := exec.CommandContext(ctx, binPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp error: %v, stderr: %s", err, stderr.String())
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return "", errors.New("empty output from yt-dlp")
	}

	lines := strings.Split(out, "\n")
	return strings.TrimSpace(lines[0]), nil
}

func EnsureYTDLP(ctx context.Context) (string, error) {
	path, err := exec.LookPath("yt-dlp")
	if err == nil {
		return path, nil // Found in system PATH
	}

	// Not found, check local config dir
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("yt-dlp not found and cannot determine home dir to download")
	}

	dir := filepath.Join(home, ".config", "soundcloud-radio")
	os.MkdirAll(dir, 0755)
	
	binName := "yt-dlp"
	if runtime.GOOS == "windows" {
		binName = "yt-dlp.exe"
	} else if runtime.GOOS == "darwin" {
		binName = "yt-dlp_macos"
	}
	
	localPath := filepath.Join(dir, "yt-dlp")
	if runtime.GOOS == "windows" {
		localPath += ".exe"
	}

	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil // Already downloaded
	}

	// Download it
	dlURL := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/" + binName
	req, err := http.NewRequestWithContext(ctx, "GET", dlURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download yt-dlp: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download yt-dlp: HTTP %d", resp.StatusCode)
	}

	out, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(localPath) // cleanup partial
		return "", err
	}
	
	// Ensure it's executable
	os.Chmod(localPath, 0755)

	return localPath, nil
}
