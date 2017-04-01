package wikirace

import (
	"bytes"
	"fmt"
	"time"
)

// Search represents a search description.
type Search struct {
	// ResultEndpoint string
	WikiEndpoint string
	Start        string
	End          string
	History      []string
	StartedAt    time.Time
	AchievedAt   time.Time
}

func (s Search) String() string {
	var b bytes.Buffer

	fmt.Fprintln(&b, "\033[00;1mResult:\033[00m")

	for _, entry := range s.History {
		fmt.Fprintf(&b, "  %v\n", entry)
	}

	ellapsed := s.AchievedAt.Sub(s.StartedAt)
	fmt.Fprintf(&b, "\033[92mFound in %.2f s\033[00m\n", ellapsed.Seconds())
	fmt.Fprintf(&b, "\033[95m%v hops\033[00m", len(s.History))
	return b.String()
}
