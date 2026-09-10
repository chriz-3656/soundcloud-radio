package soundcloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"time"
)

var (
	bundleRegex   = regexp.MustCompile(`https://a-v2\.sndcdn\.com/assets/[^"]+\.js`)
	clientIdRegex = regexp.MustCompile(`"client_id":"([a-zA-Z0-9]{32})"|client_id:"([a-zA-Z0-9]{32})"|client_id=([a-zA-Z0-9]{32})`)
)

func FetchClientID(ctx context.Context) (string, error) {
	if envID := os.Getenv("SOUNDCLOUD_CLIENT_ID"); envID != "" {
		return envID, nil
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", "https://soundcloud.com", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status fetching home: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	bundles := bundleRegex.FindAllString(string(body), -1)
	if len(bundles) == 0 {
		return "", fmt.Errorf("no JS bundles found on homepage")
	}
	
	for _, bundleURL := range bundles {
		req, err := http.NewRequestWithContext(ctx, "GET", bundleURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		content, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		
		matches := clientIdRegex.FindStringSubmatch(string(content))
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				return matches[i], nil
			}
		}
	}
	
	return "", fmt.Errorf("client_id not found in bundles")
}
