package provider

import (
	"context"
	"soundcloud-radio/internal/models"
)

type Provider interface {
	SearchTracks(ctx context.Context, q string, limit int) ([]models.Track, error)
	GetRelatedTracks(ctx context.Context, trackID string, limit int) ([]models.Track, error)
	GetHomeFeed(ctx context.Context) ([]models.Track, error)
	GetName() string
}
