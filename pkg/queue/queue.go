package queue

import (
	"sync"
)

type Queue struct {
	items []interface{}
	mu    sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]interface{}, 0),
	}
}

func (q *Queue) Push(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue) Pop() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) Remove(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, v := range q.items {
		if v == item {
			q.items = append(q.items[:i], q.items[i+1:]...)
			break
		}
	}
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
