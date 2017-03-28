package wikirace

import (
	"testing"
	"time"
)

func TestResultManagerPathFound(t *testing.T) {
	initialSearch := Search{
		Start:     "Maxence",
		End:       "Segment",
		StartedAt: time.Now(),
	}
	jobs := NewJobsStatus()
	queue := MakeSearchQueue()
	manager := NewResultManager(initialSearch, jobs, queue)

	manager.Emmit(Search{
		Start: "Segment",
		End:   "Segment",
		History: []string{
			"Maxence",
			"Segment",
		},
		StartedAt: initialSearch.StartedAt,
	})
	manager.Listen()
}

func TestResultManagerPathNotFound(t *testing.T) {
	initialSearch := Search{
		Start:     "Maxence",
		End:       "Segment",
		StartedAt: time.Now(),
	}
	jobs := NewJobsStatus()
	queue := MakeSearchQueue()
	manager := NewResultManager(initialSearch, jobs, queue)

	manager.Emmit(initialSearch)
	manager.Listen()
}

func TestResultManagerNewSearch(t *testing.T) {
	initialSearch := Search{
		Start:     "Maxence",
		End:       "Segment",
		StartedAt: time.Now(),
	}
	jobs := NewJobsStatus()
	queue := MakeSearchQueue()
	manager := NewResultManager(initialSearch, jobs, queue)

	jobs.Inc()
	manager.Emmit(initialSearch)
	jobs.Inc()
	manager.Emmit(initialSearch)

	go func() {
		time.Sleep(time.Second * 1)
		manager.Emmit(Search{
			Start:     "Segment",
			End:       "Segment",
			StartedAt: initialSearch.StartedAt,
		})
	}()

	manager.Listen()
}
