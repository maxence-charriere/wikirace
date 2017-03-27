package main

import (
	"net/http"
	"testing"

	"github.com/maxence-charriere/wikirace"
)

func TestParseHTML(t *testing.T) {
	res, err := http.Get("https://en.wikipedia.org/wiki/Maxence")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	s := wikirace.Search{
		End: "Maxence_Perrin",
	}

	h := func(s wikirace.Search) {
		t.Log(s.Start)
	}

	if err = parseHTML(res.Body, s, h); err != nil {
		t.Error(err)
	}
}

func TestOnLink(t *testing.T) {
	s := wikirace.Search{
		ResultEndpoint: "https://google.com",
	}
	onLink(s)
}
