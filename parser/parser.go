package parser

import (
	"io"
	"path"
	"strings"

	"github.com/maxence-charriere/wikirace"
	"golang.org/x/net/html"
)

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

			if strings.Contains(target, ":") {
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
