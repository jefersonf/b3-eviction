package queue

import "sync"

type MemoryQueue struct {
	mu    sync.Mutex
	items []string
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		items: make([]string, 0),
	}
}

func (q *MemoryQueue) Enqueue(item string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *MemoryQueue) Dequeue() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return "", false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
