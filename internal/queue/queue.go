package queue

import "soundcloud-radio/internal/models"

import (
	"sync"
	"math/rand"
	
	"soundcloud-radio/internal/soundcloud"
)

type Queue struct {
	mu           sync.Mutex
	items        []models.Track
	history      []models.Track
	historyLimit int
	playedIds    map[int64]bool
}

func NewQueue(historyLimit int) *Queue {
	if historyLimit <= 0 {
		historyLimit = 5000
	}
	return &Queue{
		items:        make([]models.Track, 0),
		history:      make([]models.Track, 0),
		historyLimit: historyLimit,
		playedIds:    make(map[int64]bool),
	}
}

func (q *Queue) Add(t models.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, t)
}

func (q *Queue) AddIfNotPlayed(t models.Track) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.playedIds[t.ID] {
		return false
	}
	for _, item := range q.items {
		if item.ID == t.ID {
			return false
		}
	}
	q.items = append(q.items, t)
	return true
}

func (q *Queue) Pop() (models.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return models.Track{}, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	
	q.playedIds[item.ID] = true
	q.history = append(q.history, item)
	if len(q.history) > q.historyLimit {
		q.history = q.history[1:]
	}
	
	return item, true
}

func (q *Queue) Peek() (models.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return models.Track{}, false
	}
	return q.items[0], true
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = make([]models.Track, 0)
}

func (q *Queue) Items() []models.Track {
	q.mu.Lock()
	defer q.mu.Unlock()
	itemsCopy := make([]models.Track, len(q.items))
	copy(itemsCopy, q.items)
	return itemsCopy
}

func (q *Queue) LastPlayed() (models.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.history) == 0 {
		return models.Track{}, false
	}
	return q.history[len(q.history)-1], true
}

func (q *Queue) Remove(index int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if index >= 0 && index < len(q.items) {
		q.items = append(q.items[:index], q.items[index+1:]...)
	}
}

func (q *Queue) MoveUp(index int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if index > 0 && index < len(q.items) {
		q.items[index], q.items[index-1] = q.items[index-1], q.items[index]
	}
}

func (q *Queue) MoveDown(index int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if index >= 0 && index < len(q.items)-1 {
		q.items[index], q.items[index+1] = q.items[index+1], q.items[index]
	}
}

func (q *Queue) PlayNext(track models.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append([]models.Track{track}, q.items...)
}

func (q *Queue) Shuffle() {
	q.mu.Lock()
	defer q.mu.Unlock()
	rand.Shuffle(len(q.items), func(i, j int) {
		q.items[i], q.items[j] = q.items[j], q.items[i]
	})
}
