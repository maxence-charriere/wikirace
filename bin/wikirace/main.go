package main

import (
	"log"
	"sync"
	"time"

	"github.com/maxence-charriere/wikirace"
)

var (
	searchCache = map[string]bool{}
	mtx         sync.Mutex
	searchCount uint
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Println(err)
		return
	}

	search := wikirace.Search{
		WikiEndpoint: cfg.Wiki,
		Start:        cfg.Start,
		End:          cfg.End,
		StartedAt:    time.Now(),
	}
	status := wikirace.NewJobsStatus()
	queue := wikirace.MakeSearchQueue()
	res := wikirace.NewResultManager(search, status, queue)
	jobPool := wikirace.NewJobPool(8, status, res, queue)

	queue.Enqueue(search)

	go jobPool.Start()
	defer jobPool.Stop()

	res.Listen()
}
