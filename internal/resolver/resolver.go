package resolver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"soundcloud-radio/internal/soundcloud"
)

type ResolvedStream struct {
	URL string
}

type Resolver interface {
	Resolve(ctx context.Context, track soundcloud.Track, cookieMode string) (*ResolvedStream, error)
}

type YTDLPResolver struct{}

func NewYTDLPResolver() *YTDLPResolver {
	return &YTDLPResolver{}
}

func (r *YTDLPResolver) Resolve(ctx context.Context, track soundcloud.Track, cookieMode string) (*ResolvedStream, error) {
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
	args := []string{"-g", "-f", "bestaudio", "--no-warnings", "--no-playlist"}
	if cookies != "" {
		args = append(args, "--cookies-from-browser", cookies)
	}
	args = append(args, trackURL)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
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
