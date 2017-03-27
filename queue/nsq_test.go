package queue

import (
	"encoding/json"
	"testing"

	"github.com/maxence-charriere/wikirace"
)

func TestNsqEnqueuer(t *testing.T) {
	deq, err := NewNsqEnqueuer(":4150")
	if err != nil {
		t.Fatal(err)
	}
	defer deq.Close()

	s := wikirace.Search{
		Start: "Maxence",
		End:   "Rome",
	}

	if err = deq.Enqueue(s); err != nil {
		t.Error(err)
	}
}

func TestNsqDequeuer(t *testing.T) {
	deq, err := NewNsqDequeuer(":4150")
	if err != nil {
		t.Fatal(err)
	}
	defer deq.Close()

	var res wikirace.Search

	h := func(body []byte) error {
		json.Unmarshal(body, &res)
		deq.StopDequeue()
		return nil
	}

	if err := deq.StartDequeue(h); err != nil {
		t.Fatal(err)
	}

	if res.Start != "Maxence" {
		t.Error("res.Start should me Maxence:", res.Start)
	}
	if res.End != "Rome" {
		t.Error("res.End should me Rome:", res.End)
	}
}
