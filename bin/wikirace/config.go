package main

import (
	"flag"
	"strings"

	"github.com/pkg/errors"
)

type conf struct {
	Wiki  string
	Start string
	End   string
	Bind  string
	Queue string
}

func getConfig() (cfg conf, err error) {
	flag.StringVar(&cfg.Wiki, "wiki", "https://en.wikipedia.org", "The wiki root where the race take place.")
	flag.StringVar(&cfg.Start, "start", "Mike Tyson", "The start of the race.")
	flag.StringVar(&cfg.End, "end", "Segment", "The end of the race.")
	flag.StringVar(&cfg.Bind, "bind", ":8042", "The bind addr.")
	flag.StringVar(&cfg.Queue, "queue", ":4150", "The queue addr.")
	flag.Parse()

	if len(cfg.Start) == 0 {
		err = errors.New("start can't be empty")
		return
	}
	if len(cfg.End) == 0 {
		err = errors.New("end can't be empty")
		return
	}

	cfg.Start = epureTarget(cfg.Start)
	cfg.End = epureTarget(cfg.End)
	return
}

func epureTarget(t string) string {
	t = strings.TrimSpace(t)
	return strings.Replace(t, " ", "_", -1)
}
