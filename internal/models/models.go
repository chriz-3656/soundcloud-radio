package models

type Track struct {
	ID           string
	Title        string
	Artist       string
	Album        string
	Duration     int64
	PermalinkURL string
	ArtworkURL   string
	DirectMediaURL string
}
