package main

import (
	"io/ioutil"
	"testing"

	"github.com/azuretek/crawler/internal/job"
	uuid "github.com/satori/go.uuid"
)

func TestGetURLs(t *testing.T) {
	example1, err := ioutil.ReadFile("../../fixtures/example1.html")
	if err != nil {
		t.Errorf("couldn't load example 1: %v", err)
	}

	t.Run("Example 1", func(t *testing.T) {
		url := &job.URL{
			URL:   "http://example.com",
			ID:    uuid.NewV4(),
			Depth: 0,
		}

		urls, err := GetURLs(url, example1)
		if err != nil {
			t.Error(err)
		}

		for _, u := range urls {
			if u.Depth != 1 {
				t.Errorf("%v depth set incorrectly", u.URL)
			}
		}

		if urls[0].URL != "http://example.com/about.html" {
			t.Error("url 0 doesn't match expected value")
		}

		if urls[1].URL != "http://example.com/static" {
			t.Error("url 1 doesn't match expected value")
		}
	})

	example2, err := ioutil.ReadFile("../../fixtures/example2.html")
	if err != nil {
		t.Errorf("couldn't load example 2: %v", err)
	}

	t.Run("Example 2", func(t *testing.T) {
		url := &job.URL{
			URL:   "http://example.com",
			ID:    uuid.NewV4(),
			Depth: 1,
		}

		urls, err := GetURLs(url, example2)
		if err != nil {
			t.Error(err)
		}

		for _, u := range urls {
			if u.Depth != 2 {
				t.Errorf("%v depth set incorrectly", u.URL)
			}
		}

		if urls[0].URL != "http://example.com/index.html" {
			t.Error("url 0 doesn't match expected value")
		}

		if urls[1].URL != "http://support.com" {
			t.Error("url 1 doesn't match expected value")
		}

		if urls[2].URL != "http://google.com" {
			t.Error("url 2 doesn't match expected value")
		}

		if urls[3].URL != "https://support.com/example" {
			t.Error("url 3 doesn't match expected value")
		}
	})
}
