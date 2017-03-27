package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/maxence-charriere/wikirace"
	"github.com/maxence-charriere/wikirace/queue"
)

var (
	enqueuer    wikirace.Enqueuer
	mtx         sync.Mutex
	searchCache = map[string]bool{}
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Println(err)
		return
	}

	if enqueuer, err = queue.NewNsqEnqueuer(cfg.Queue); err != nil {
		log.Fatalln("connection to queue failed:", err)
	}
	defer enqueuer.Close()

	wikiURL, err := url.Parse("http://" + cfg.Bind)
	if err != nil {
		log.Fatalln("bad bind addr:", err)
	}
	wikiURL.Path = path.Join(wikiURL.Path, "results")

	s := wikirace.Search{
		ResultEndpoint: wikiURL.String(),
		WikiEnpoint:    cfg.Wiki,
		Start:          cfg.Start,
		End:            cfg.End,
		History:        []string{cfg.Start},
		StartedAt:      time.Now(),
	}

	searchCache[s.Start] = true
	go func() {
		if err = enqueuer.Enqueue(s); err != nil {
			log.Fatalln("unable to queue a search:", err)
		}
	}()

	http.HandleFunc("/results", onResult)
	if err := http.ListenAndServe(cfg.Bind, nil); err != nil {
		log.Fatalln("server error:", err)
	}
}

func onResult(res http.ResponseWriter, req *http.Request) {
	var s wikirace.Search
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&s); err != nil {
		log.Println("decode result failed:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.History = append(s.History, s.Start)

	mtx.Lock()
	defer mtx.Unlock()

	// log.Println(s.Start)

	if strings.ToLower(s.Start) == strings.ToLower(s.End) {
		printSearchResult(s)
		os.Exit(0)
	}

	if _, ok := searchCache[s.Start]; ok {
		log.Println("~> Duplicate", s.Start)
		res.WriteHeader(http.StatusOK)
		return
	}

	log.Println("~>REAL!!!!!!", s.Start)
	searchCache[s.Start] = true
	if err := enqueuer.Enqueue(s); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func printSearchResult(s wikirace.Search) {
	ellapsed := time.Now().Sub(s.StartedAt)

	for _, h := range s.History {
		fmt.Println(h)
	}

	fmt.Println("found in", ellapsed.Nanoseconds(), "ns")

}
