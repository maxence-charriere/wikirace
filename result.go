package wikirace

import (
	"fmt"
	"strings"
	"time"
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
	jobQueue     Queuer
	resQueue     *SearchQueue
}

func NewResultManager(initialSeach Search, s JobsStatuer, q Queuer) *ResultManager {
	return &ResultManager{
		intialSearch: initialSeach,
		jobs:         s,
		jobQueue:     q,
		resQueue:     NewSearchQueue(),
	}
}

func (m *ResultManager) Emmit(res Search) {
	m.resQueue.Enqueue(res)
}

func (m *ResultManager) Listen() {

	for {
		res, ok := m.resQueue.Dequeue()
		if !ok {
			time.Sleep(time.Millisecond)
			continue
		}

		if len(res.Start) == len(res.End) &&
			strings.ToLower(res.Start) == strings.ToLower(res.End) {
			historyEntry := strings.Replace(res.Start, "_", " ", -1)

			res.History = append(res.History, historyEntry)
			res.AchievedAt = time.Now()

			fmt.Println(res)
			fmt.Printf("\033[91m%v links processed\033[00m\n\n", m.jobs.Total())
			break
		}

		// No path found.
		if m.jobQueue.Len() == 0 && m.jobs.Count() == 0 {
			fmt.Println("\033[00;1mResult:\033[00m")
			fmt.Printf("  \033[91mno path found for %v to %v\033[00m\n\n",
				m.intialSearch.Start,
				m.intialSearch.End,
			)
			break
		}

		// New search to queue.
		m.jobQueue.Enqueue(res)
	}
	return
}
