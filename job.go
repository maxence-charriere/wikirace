package wikirace

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// JobsStatuer is the interface to get or modify the overall job status.
type JobsStatuer interface {
	Inc()
	Dec()
	Count() int
	Total() int
	IsHandled(s Search) bool
	SetHandled(s Search)
}

// JobsStatus represent the overall state of the jobs currently processed.
// Thread safe.
type JobsStatus struct {
	mutex       sync.Mutex
	count       int
	handledJobs map[string]bool
}

func NewJobsStatus() *JobsStatus {
	return &JobsStatus{
		handledJobs: map[string]bool{},
	}
}

func (s *JobsStatus) Inc() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count++
}

func (s *JobsStatus) Dec() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count--
	if s.count < 0 {
		s.count = 0
	}
}

func (s *JobsStatus) Count() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.count
}

func (s *JobsStatus) Total() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.handledJobs)
}

func (s *JobsStatus) IsHandled(search Search) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.handledJobs[search.Start]
	return ok
}

func (s *JobsStatus) SetHandled(search Search) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.handledJobs[search.Start] = true
}

type JobPooler interface {
	Start()
	Stop()
}

type JobPool struct {
	maxRoutines int
	jobs        JobsStatuer
	res         ResultEmitter
	queue       Dequeuer
	stopChan    chan interface{}
}

func NewJobPool(maxRoutines int, s JobsStatuer, e ResultEmitter, q Dequeuer) *JobPool {
	return &JobPool{
		maxRoutines: maxRoutines,
		jobs:        s,
		res:         e,
		queue:       q,
		stopChan:    make(chan interface{}, 1),
	}
}

func (p *JobPool) Start() {
	for {
		select {
		case <-p.stopChan:
			return

		default:
			p.jobs.Inc()
			if search, ok := p.queue.Dequeue(); ok {
				p.doJob(search)
			} else {
				time.Sleep(time.Millisecond)
			}
			p.jobs.Dec()
		}
	}
}

func (p *JobPool) Stop() {
	p.stopChan <- nil
}

func (p *JobPool) doJob(s Search) {
	if p.jobs.IsHandled(s) {
		return
	}
	p.jobs.SetHandled(s)

	fmt.Println("searching in", s.Start)
	historyEntry := strings.Replace(s.Start, "_", " ", -1)
	s.History = append(s.History, historyEntry)

	URL, err := url.Parse(s.WikiEndpoint)
	if err != nil {
		fmt.Printf("\033[91m%v ~> KO: %v\n\033[00m", s.Start, err)
		return
	}
	URL.Path = path.Join(URL.Path, "wiki", s.Start)

	res, err := http.Get(URL.String())
	if err != nil {
		fmt.Printf("\033[91m%v ~> KO: %v\n\033[00m", s.Start, err)
		return
	}
	defer res.Body.Close()

	links, err := p.parseHTML(res.Body)
	if err != nil {
		fmt.Printf("\033[91m%v ~> KO: %v\n\033[00m", s.Start, err)
		return
	}

	for _, link := range links {
		newSearch := s
		newSearch.Start = link
		newSearch.History = make([]string, len(s.History))
		copy(newSearch.History, s.History)
		p.res.Emmit(newSearch)
	}
}

func (p *JobPool) parseHTML(body io.Reader) (links []string, err error) {
	titleScope := false

	z := html.NewTokenizer(body)
	for {
		switch token := z.Next(); token {
		case html.ErrorToken:
			if err = z.Err(); err != io.EOF {
				return
			}
			err = nil
			return

		case html.StartTagToken:
			switch tag, _ := z.TagName(); string(tag) {
			case "title":
				titleScope = true

			case "a":
				if l, ok := parseLink(z); ok {
					links = append(links, l)
				}
			}

		case html.TextToken:
			if !titleScope {
				continue
			}
			if l, ok := parseTitle(z); ok {
				links = append(links, l)
			}
			titleScope = false
		}
	}
}

func parseTitle(z *html.Tokenizer) (link string, ok bool) {
	title := string(z.Text())
	splited := strings.Split(title, " -")

	if len(splited) == 0 {
		return
	}

	link = strings.TrimSpace(splited[0])
	ok = true
	return
}

func parseLink(z *html.Tokenizer) (link string, ok bool) {
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
		return
	}

	// Exlude non wiki pages.
	dir, name := path.Split(href)
	if dir != "/wiki/" {
		return
	}

	// Exclude special pages.
	if strings.Contains(name, ":") {
		return
	}

	link = name
	ok = true
	return
}
