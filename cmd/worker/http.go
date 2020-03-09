package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"
)

// GetBody returns the body of the requested URL
func GetBody(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetHref returns the value of the html href attribute
func GetHref(attr []html.Attribute) (string, error) {
	for _, a := range attr {
		if a.Key == "href" {
			if a.Val == "" {
				return "", fmt.Errorf("href value empty")
			}
			return a.Val, nil
		}
	}

	return "", fmt.Errorf("no href found")
}
