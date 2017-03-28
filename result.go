package wikirace

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ResultEmitter is the interface to emit a search result.
// To be used when a new link is found.
type ResultEmitter interface {
	Emmit(res Search)
}

// ResultListener is the interface that launch the loop that wait for search
// results.
type ResultListener interface {
	Listen() (err error)
}

// ResultManager implements ResultEmitter and ResultListener. It is in charge
// to decide what to do when a search operation is completed.
type ResultManager struct {
	intialSearch Search
	jobs         JobsStatuer
	queue        Queuer
	resChan      chan Search
}

func NewResultManager(initialSeach Search, s JobsStatuer, q Queuer) *ResultManager {
	return &ResultManager{
		intialSearch: initialSeach,
		jobs:         s,
		queue:        q,
		resChan:      make(chan Search, 256),
	}
}

func (m *ResultManager) Emmit(res Search) {
	m.resChan <- res
}

func (m *ResultManager) Listen() (err error) {
	for res := range m.resChan {
		// Complete path found.
		if len(res.Start) == len(res.End) &&
			strings.ToLower(res.Start) == strings.ToLower(res.End) {
			res.AchievedAt = time.Now()
			fmt.Println(res)
			break
		}

		// No path found.
		if m.queue.Len() == 0 && m.jobs.Count() == 0 {
			err = errors.Errorf("no path found for %v to %v",
				m.intialSearch.Start,
				m.intialSearch.End)
			fmt.Println(err)
			break
		}

		// New search to queue.
		if m.jobs.IsHandled(res) {
			continue
		}
		m.queue.Enqueue(res)
		m.jobs.SetHandled(res)
	}
	return
}
