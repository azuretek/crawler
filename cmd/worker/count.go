package main

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Host contains information about a specific hostname
type Host struct {
	Name  string
	Count int64
}

// Increment adds 1 to count
func (h *Host) Increment() int64 {
	h.Count = h.Count + 1
	return h.Count
}

// CountHostNames returns a map of Host objects
func CountHostNames(body []byte) (map[string]*Host, error) {
	hostnames := map[string]*Host{}
	t := html.NewTokenizer(strings.NewReader(string(body)))
	for {
		tt := t.Next()

		switch tt {

		// Return gracefully if we finish, otherwise return an error
		case html.ErrorToken:
			if t.Err() == io.EOF {
				return hostnames, nil
			}
			return hostnames, t.Err()

		// Count hostnames from anchor tags
		case html.StartTagToken:
			token := t.Token()
			if token.Data == "a" {
				href, err := GetHref(token.Attr)
				if err != nil {
					// If couldn't find href skip anchor
					continue
				}

				u, err := url.Parse(href)
				if err != nil {
					return hostnames, err
				}

				// might return host:port and we just want host
				host := strings.Split(u.Host, ":")[0]

				// if host is empty skip adding/incrementing a host
				// this is usually caused by relative paths in the href
				if host == "" {
					continue
				}

				// if host already exists increment, otherwise add new host
				if h, ok := hostnames[host]; ok {
					h.Increment()
				} else {
					hostnames[host] = &Host{
						Name:  host,
						Count: 1,
					}
				}
			}
		}
	}
}
