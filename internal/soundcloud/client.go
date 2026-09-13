package soundcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"soundcloud-radio/internal/models"
)

type Client struct {
	BaseURL  string
	ClientID string
	HTTP     *http.Client
}

func NewClient(baseURL, clientID string) *Client {
	if baseURL == "" {
		baseURL = "https://api-v2.soundcloud.com"
	}
	return &Client{
		BaseURL:  baseURL,
		ClientID: clientID,
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) GetName() string {
	return "SoundCloud"
}

func (c *Client) doRequest(ctx context.Context, path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("client_id", c.ClientID)
	
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()
	
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}

func (c *Client) SearchTracks(ctx context.Context, q string, limit int) ([]models.Track, error) {
	query := url.Values{}
	query.Set("q", q)
	query.Set("limit", fmt.Sprint(limit))
	
	body, err := c.doRequest(ctx, "/search/tracks", query)
	if err != nil {
		return nil, err
	}
	
	var sr SearchResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}
	
	var tracks []models.Track
	for _, t := range sr.Collection {
		tracks = append(tracks, models.Track{
			ID:           fmt.Sprint(t.ID),
			Title:        t.Title,
			Artist:       t.User.Username,
			Duration:     t.Duration,
			PermalinkURL: t.PermalinkURL,
			ArtworkURL:   t.ArtworkURL,
		})
	}
	
	return tracks, nil
}

func (c *Client) GetRelatedTracks(ctx context.Context, trackID string, limit int) ([]models.Track, error) {
	query := url.Values{}
	query.Set("limit", fmt.Sprint(limit))
	
	path := fmt.Sprintf("/tracks/%s/related", trackID)
	body, err := c.doRequest(ctx, path, query)
	if err != nil {
		return nil, err
	}
	
	var sr SearchResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}
	
	var tracks []models.Track
	for _, t := range sr.Collection {
		tracks = append(tracks, models.Track{
			ID:           fmt.Sprint(t.ID),
			Title:        t.Title,
			Artist:       t.User.Username,
			Duration:     t.Duration,
			PermalinkURL: t.PermalinkURL,
			ArtworkURL:   t.ArtworkURL,
		})
	}
	
	return tracks, nil
}
