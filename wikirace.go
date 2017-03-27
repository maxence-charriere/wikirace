package wikirace

import "time"

type Search struct {
	ResultEndpoint string
	WikiEnpoint    string
	Start          string
	End            string
	History        []string
	StartedAt      time.Time
}

type DequeueHandler func(body []byte) error

type Enqueuer interface {
	Enqueue(s Search) error
	Close()
}

type Dequeuer interface {
	StartDequeue(h DequeueHandler) error
	StopDequeue()
	Close()
}

type LinkParser interface {
}
