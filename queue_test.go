package wikirace

import "testing"

func TestSearchQueue(t *testing.T) {
	q := MakeSearchQueue(42)
	q.Enqueue(Search{
		Start: "Maxence",
	})

	if l := len(q); l != 1 {
		t.Error("l should be 1:", l)
	}

	s, ok := q.Dequeue()
	if !ok {
		t.Fatal("ok should be true")
	}
	if s.Start != "Maxence" {
		t.Error("s.Start should be Maxence:", s.Start)
	}

	if l := q.Len(); l != 0 {
		t.Error("l should be 0:", l)
	}

	if s, ok = q.Dequeue(); ok {
		t.Fatal("ok should be false")
	}
}
