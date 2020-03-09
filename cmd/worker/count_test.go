package main

import (
	"io/ioutil"
	"testing"
)

func TestIncrement(t *testing.T) {
	host := &Host{
		Name:  "foo.com",
		Count: 10,
	}

	if host.Increment() != 11 {
		t.Error("Increment didn't add 1")
	}
}

func TestCountHostNames(t *testing.T) {
	example1, err := ioutil.ReadFile("../../fixtures/example1.html")
	if err != nil {
		t.Errorf("couldn't load example 1: %v", err)
	}

	t.Run("Example 1", func(t *testing.T) {
		h, err := CountHostNames(example1)
		if err != nil {
			t.Error(err)
		}

		if _, ok := h["example.com"]; !ok {
			t.Error("Couldn't find example.com hostname")
		}

		if h["example.com"].Count != 1 {
			t.Error("example.com count not equal to 1")
		}
	})

	example2, err := ioutil.ReadFile("../../fixtures/example2.html")
	if err != nil {
		t.Errorf("couldn't load example 2: %v", err)
	}

	t.Run("Example 2", func(t *testing.T) {
		h, err := CountHostNames(example2)
		if err != nil {
			t.Error(err)
		}

		if _, ok := h["support.com"]; !ok {
			t.Error("Couldn't find support.com hostname")
		}

		if h["support.com"].Count != 2 {
			t.Error("support.com count not equal to 2")
		}

		if _, ok := h["google.com"]; !ok {
			t.Error("Couldn't find google.com hostname")
		}

		if h["google.com"].Count != 1 {
			t.Error("google.com count not equal to 2")
		}
	})
}
