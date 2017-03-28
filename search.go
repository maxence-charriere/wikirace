package wikirace

import "time"

type Search struct {
	ResultEndpoint string
	WikiEndpoint   string
	Start          string
	End            string
	History        []string
	StartedAt      time.Time
	AchievedAt     time.Time
}
