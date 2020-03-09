package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/azuretek/crawler/internal/job"
	"golang.org/x/net/html"
)

// GetURLs returns a list of URL objects
func GetURLs(u *job.URL, body []byte) ([]*job.URL, error) {
	urls := []*job.URL{}
	t := html.NewTokenizer(strings.NewReader(string(body)))
	for {
		tt := t.Next()

		switch tt {

		// return gracefully if we finish, otherwise return an error
		case html.ErrorToken:
			if t.Err() == io.EOF {
				return urls, nil
			}
			return urls, t.Err()

		// get urls from anchor tags
		case html.StartTagToken:
			token := t.Token()
			if token.Data == "a" {
				href, err := GetHref(token.Attr)
				if err != nil {
					// If couldn't find href skip anchor
					continue
				}

				url := href
				if strings.HasPrefix(href, "?") || !strings.HasPrefix(href, "http") {
					url = fmt.Sprintf("%v/%v", u.URL, href)
				}
				if strings.HasPrefix(href, "/") {
					url = fmt.Sprintf("%v%v", u.URL, href)
				}

				if !urlExists(urls, url) {
					urls = append(urls, &job.URL{
						URL:   url,
						ID:    u.ID,
						Depth: u.Depth + 1,
					})
				}
			}
		}
	}
}

func urlExists(urls []*job.URL, url string) bool {
	for _, u := range urls {
		if u.URL == url {
			return true
		}
	}
	return false
}
