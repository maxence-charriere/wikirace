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
	queue := MakeSearchQueue(42)
	manager := NewResultManager(initialSearch, jobs, queue)

	manager.Emmit(Search{
		Start:     "Segment",
		End:       "Segment",
		StartedAt: initialSearch.StartedAt,
	})

	if err := manager.Listen(); err != nil {
		t.Error(err)
	}
}

func TestResultManagerPathNotFound(t *testing.T) {
	initialSearch := Search{
		Start:     "Maxence",
		End:       "Segment",
		StartedAt: time.Now(),
	}
	jobs := NewJobsStatus()
	queue := MakeSearchQueue(42)
	manager := NewResultManager(initialSearch, jobs, queue)

	manager.Emmit(initialSearch)

	if err := manager.Listen(); err == nil {
		t.Error("err should not be nil")
	}
}

func TestResultManagerNewSearch(t *testing.T) {
	initialSearch := Search{
		Start:     "Maxence",
		End:       "Segment",
		StartedAt: time.Now(),
	}
	jobs := NewJobsStatus()
	queue := MakeSearchQueue(42)
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

	if err := manager.Listen(); err != nil {
		t.Error(err)
	}
}
