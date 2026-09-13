package soundcloud

type SearchResponse struct {
	Collection []TrackJSON `json:"collection"`
}

type TrackJSON struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Duration     int64  `json:"duration"`
	PermalinkURL string `json:"permalink_url"`
	ArtworkURL   string `json:"artwork_url"`
	User         struct {
		Username string `json:"username"`
	} `json:"user"`
}
