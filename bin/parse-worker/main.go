package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/maxence-charriere/wikirace"
	"github.com/maxence-charriere/wikirace/queue"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type searchCompetionHandler func(s wikirace.Search) error

func main() {
	var queueEndpoint string
	flag.StringVar(&queueEndpoint, "queueURL", ":4150", "Addr of the queue to connect.")
	flag.Parse()

	deq, err := queue.NewNsqDequeuer(queueEndpoint)
	if err != nil {
		log.Fatalln("can't connect to the queue:", err)
	}
	defer deq.Close()

	if err = deq.StartDequeue(onDequeue); err != nil {
		log.Fatalln("can't connect to the queue:", err)
	}
}

func onDequeue(body []byte) error {
	var s wikirace.Search
	if err := json.Unmarshal(body, &s); err != nil {
		log.Panicf("can't unmarshal %s: %v", body, err)
	}

	URL, err := url.Parse(s.WikiEnpoint)
	if err != nil {
		log.Println("url parsing failed:", err)
		return err
	}
	URL.Path = path.Join(URL.Path, "wiki", s.Start)

	c := &http.Client{}
	res, err := c.Get(URL.String())
	if err != nil {
		log.Println("get request failed:", err)
		return err
	}
	defer res.Body.Close()

	if err = parseHTML(res.Body, s, onLink); err != nil {

	}
	return err
}

func parseHTML(body io.Reader, search wikirace.Search, h searchCompetionHandler) error {
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

			wikiLink := ""
			for {
				key, val, more := z.TagAttr()
				if string(key) == "href" {
					wikiLink = string(val)
					break
				}
				if !more {
					break
				}
			}

			if len(wikiLink) == 0 {
				continue
			}

			wikiLink = path.Clean(wikiLink)
			dir, target := path.Split(wikiLink)
			if dir != "/wiki/" {
				continue
			}

			search.Start = target
			if err := h(search); err != nil {
				return err
			}

			if search.End == target {
				return nil
			}
		}
	}
}

func onLink(s wikirace.Search) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	if err := enc.Encode(s); err != nil {
		return errors.Wrapf(err, "result encoding for %s", s.Start)
	}

	c := &http.Client{}
	res, err := c.Post(s.ResultEndpoint, "application/json", &b)
	if err != nil {
		log.Printf("can't send result for %v: %v", s.Start, err)
		return nil
	}
	res.Body.Close()
	return nil
}
