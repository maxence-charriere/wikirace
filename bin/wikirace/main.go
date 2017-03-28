package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/davecgh/go-spew/spew"
)

var (
	searchCache = map[string]bool{}
	mtx         sync.Mutex
	searchCount uint
)

type Search struct {
	WikiEnpoint string
	Start       string
	End         string
	History     []string
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Println(err)
		return
	}

	searchChan := make(chan Search, 1000)
	resChan := make(chan Search)

	go startSearchWorker(searchChan, resChan)

	searchChan <- Search{
		WikiEnpoint: cfg.Wiki,
		Start:       cfg.Start,
		End:         cfg.End,
	}

	res := <-resChan
	spew.Dump(res)
}

func startSearchWorker(searchChan chan Search, resChan chan<- Search) {
	for s := range searchChan {
		if searchCount <= 8 {
			go searchJob(s, searchChan, resChan)
			continue
		}

		if searchCount == 0 {
			log.Fatal("no path")
		}
		time.Sleep(time.Millisecond)
	}
}

func searchJob(s Search, searchChan chan<- Search, resChan chan<- Search) {
	mtx.Lock()
	searchCount++
	mtx.Unlock()

	defer func() {
		mtx.Lock()
		searchCount--
		mtx.Unlock()
	}()

	fmt.Println("starting job for", s.Start)
	s.History = append(s.History, strings.Replace(s.Start, "_", " ", -1))

	URL, err := url.Parse(s.WikiEnpoint)
	if err != nil {
		log.Println("url parsing failed:", err)
		return
	}
	URL.Path = path.Join(URL.Path, "wiki", s.Start)

	res, err := http.Get(URL.String())
	if err != nil {
		log.Println("get request failed:", err)
		return
	}
	defer res.Body.Close()

	parseHTML(res.Body, s, searchChan, resChan)
}

func parseHTML(body io.Reader, s Search, searchChan chan<- Search, resChan chan<- Search) error {
	z := html.NewTokenizer(body)

	for {
		switch token := z.Next(); token {
		case html.ErrorToken:
			if err := z.Err(); err != io.EOF {
				return err
			}
			return nil

		case html.StartTagToken:
			if tag, _ := z.TagName(); tag[0] != 'a' {
				continue
			}

			href := ""
			for {
				key, val, more := z.TagAttr()
				if string(key) == "href" {
					href = string(val)
					break
				}

				if !more {
					break
				}
			}

			if len(href) == 0 {
				continue
			}

			dir, target := path.Split(href)
			if dir != "/wiki/" {
				continue
			}

			if strings.ToLower(target) == strings.ToLower(s.End) {
				s.History = append(s.History, strings.Replace(target, "_", " ", -1))

				resChan <- s
				return nil
			}

			if strings.Contains(target, ":") {
				continue
			}

			newSearch := s
			newSearch.Start = target
			newSearch.History = make([]string, len(s.History))
			copy(newSearch.History, s.History)

			mtx.Lock()
			if _, ok := searchCache[target]; ok {
				mtx.Unlock()
				continue
			}
			searchCache[target] = true
			go func() { searchChan <- newSearch }()
			mtx.Unlock()
		}
	}
}
