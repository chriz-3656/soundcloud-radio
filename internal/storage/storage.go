package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"soundcloud-radio/internal/soundcloud"
)

type HistoryItem struct {
	Track     soundcloud.Track `json:"track"`
	PlayedAt  time.Time        `json:"played_at"`
}

type StoreData struct {
	Favorites []soundcloud.Track `json:"favorites"`
	History   []HistoryItem      `json:"history"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	Data StoreData
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.TempDir()
	}
	
	dir := filepath.Join(home, ".config", "soundcloud-radio")
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "data.json")

	s := &Store{
		path: path,
		Data: StoreData{
			Favorites: make([]soundcloud.Track, 0),
			History:   make([]HistoryItem, 0),
		},
	}

	s.Load()
	return s, nil
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.Data)
}

func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}

func (s *Store) AddFavorite(t soundcloud.Track) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, f := range s.Data.Favorites {
		if f.ID == t.ID {
			return // Already exists
		}
	}
	s.Data.Favorites = append(s.Data.Favorites, t)
	go s.Save()
}

func (s *Store) RemoveFavorite(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var updated []soundcloud.Track
	for _, f := range s.Data.Favorites {
		if f.ID != id {
			updated = append(updated, f)
		}
	}
	s.Data.Favorites = updated
	go s.Save()
}

func (s *Store) IsFavorite(id int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.Data.Favorites {
		if f.ID == id {
			return true
		}
	}
	return false
}

func (s *Store) AddHistory(t soundcloud.Track) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.Data.History = append([]HistoryItem{{Track: t, PlayedAt: time.Now()}}, s.Data.History...)
	if len(s.Data.History) > 500 {
		s.Data.History = s.Data.History[:500]
	}
	go s.Save()
}

func (s *Store) GetHistory() []HistoryItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hist := make([]HistoryItem, len(s.Data.History))
	copy(hist, s.Data.History)
	return hist
}

func (s *Store) GetFavorites() []soundcloud.Track {
	s.mu.RLock()
	defer s.mu.RUnlock()
	favs := make([]soundcloud.Track, len(s.Data.Favorites))
	copy(favs, s.Data.Favorites)
	return favs
}
