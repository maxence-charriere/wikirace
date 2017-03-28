package wikirace

import (
	"testing"
	"time"
)

func TestJobsStatus(t *testing.T) {
	jobs := NewJobsStatus()

	s := Search{
		Start: "Maxence",
	}
	jobs.SetHandled(s)
	if !jobs.IsHandled(s) {
		t.Error("s should be handled")
	}

	jobs.Inc()
	if c := jobs.Count(); c != 1 {
		t.Error("c should be 1:", c)
	}

	jobs.Dec()
	if c := jobs.Count(); c != 0 {
		t.Error("c should be 0:", c)
	}
	jobs.Dec()
	if c := jobs.Count(); c != 0 {
		t.Error("c should be 0:", c)
	}
}

func TestJobPool(t *testing.T) {
	initialSearch := Search{
		Start:        "Maxence",
		End:          "Segment",
		WikiEndpoint: "https://en.wikipedia.org",
		StartedAt:    time.Now(),
	}
	jobs := NewJobsStatus()
	queue := MakeSearchQueue()
	res := NewResultManager(initialSearch, jobs, queue)

	queue.Enqueue(initialSearch)

	pool := NewJobPool(8, jobs, res, queue)
	go func() {
		time.Sleep(time.Millisecond * 50)
		pool.Stop()
	}()
	pool.Start()
}
