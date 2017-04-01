package wikirace

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
// Implemented on the top of go channel.
type SearchQueue chan Search

func MakeSearchQueue() SearchQueue {
	return make(SearchQueue, 4096)
}

func (q SearchQueue) Len() int {
	return len(q)
}

// Enqueue enqueues s.
// Block if q is full.
func (q SearchQueue) Enqueue(s Search) error {
	q <- s
	return nil
}

// Dequeue dequeue s.
// ok will be false if the queue is empty.
func (q SearchQueue) Dequeue() (s Search, ok bool) {
	// select {
	// case s = <-q:
	// 	ok = true
	// 	return

	// default:
	// 	return
	// }

	s = <-q
	ok = true
	return
}
