package main

import (
	"testing"

	"golang.org/x/net/html"
)

func TestGetHref(t *testing.T) {
	attr := []html.Attribute{}
	attr = append(attr, html.Attribute{
		Key: "href",
		Val: "http://foo.com",
	})

	url, err := GetHref(attr)
	if err != nil {
		t.Errorf("GetHref failed: %v", err)
	}

	if url != "http://foo.com" {
		t.Error("GetHref didn't return the right href")
	}
}
