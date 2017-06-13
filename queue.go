package wikirace

import "github.com/foize/go.fifo"

// Queuer is the interface that wraps an enqueue operation.
type Queuer interface {
	Len() int
	Enqueue(s Search) error
}

// Dequeuer is the interface that wraps an dequeue operation.
type Dequeuer interface {
	Len() int
	Dequeue() (s Search, ok bool)
}

// SearchQueue is a queue used as pipeline for search operation.
type SearchQueue struct {
	queue *fifo.Queue
}

func NewSearchQueue() *SearchQueue {
	return &SearchQueue{
		queue: fifo.NewQueue(),
	}
}

func (q *SearchQueue) Len() int {
	return q.queue.Len()
}

// Enqueue enqueues s.
func (q *SearchQueue) Enqueue(s Search) error {
	q.queue.Add(s)
	return nil
}

// Dequeue dequeue s.
// ok will be false if the queue is empty.
func (q *SearchQueue) Dequeue() (s Search, ok bool) {
	if q.queue.Len() == 0 {
		return
	}

	s = q.queue.Next().(Search)
	ok = true
	return
}
