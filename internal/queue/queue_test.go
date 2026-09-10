package queue

import (
	"sync"
	"testing"
	"soundcloud-radio/internal/soundcloud"
)

func TestQueueConcurrent(t *testing.T) {
	q := NewQueue(10)
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			q.Add(soundcloud.Track{ID: int64(id)})
			q.Pop()
			q.Len()
		}(i)
	}
	wg.Wait()
}
