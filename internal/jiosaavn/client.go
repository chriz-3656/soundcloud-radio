package jiosaavn

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
	BaseURL string
	HTTP    *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://www.jiosaavn.com/api.php",
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) GetName() string {
	return "JioSaavn"
}

func (c *Client) doRequest(ctx context.Context, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("_format", "json")
	query.Set("_marker", "0")

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0")
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
	query.Set("__call", "search.getResults")
	query.Set("q", q)
	query.Set("p", "1")
	query.Set("n", fmt.Sprint(limit))
	query.Set("ctx", "web6dot0")

	body, err := c.doRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	var sr struct {
		Results []struct {
			ID             string `json:"id"`
			Song           string `json:"song"`
			PrimaryArtists string `json:"primary_artists"`
			PermaURL       string `json:"perma_url"`
			Image          string `json:"image"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}

	var tracks []models.Track
	for _, t := range sr.Results {
		tracks = append(tracks, models.Track{
			ID:           t.ID,
			Title:        t.Song,
			Artist:       t.PrimaryArtists,
			PermalinkURL: t.PermaURL,
			ArtworkURL:   t.Image,
			Duration:     0, // Need extra parsing for duration, keeping it 0 for now
		})
	}
	return tracks, nil
}

func (c *Client) GetRelatedTracks(ctx context.Context, trackID string, limit int) ([]models.Track, error) {
	query := url.Values{}
	query.Set("__call", "reco.getreco")
	query.Set("pid", trackID)
	query.Set("ctx", "android") // Must use android for reco

	body, err := c.doRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	// Dynamic json object where key is the trackID or array
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var tracks []models.Track
	// Depending on response, we might get an array directly or a map. Let's try map first
	for _, v := range raw {
		if list, ok := v.([]interface{}); ok {
			for _, item := range list {
				if tmap, ok := item.(map[string]interface{}); ok {
					id, _ := tmap["id"].(string)
					song, _ := tmap["song"].(string)
					artist, _ := tmap["primary_artists"].(string)
					purl, _ := tmap["perma_url"].(string)
					img, _ := tmap["image"].(string)
					
					tracks = append(tracks, models.Track{
						ID:           id,
						Title:        song,
						Artist:       artist,
						PermalinkURL: purl,
						ArtworkURL:   img,
					})
				}
			}
			break
		}
	}
	
	if len(tracks) > limit {
		tracks = tracks[:limit]
	}
	
	return tracks, nil
}
