package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLinkUnmarshal(t *testing.T) {
	link := Link{}
	url := "http://example.com"
	str := fmt.Sprintf(`{"url": "%s"}`, url)
	err := json.Unmarshal([]byte(str), &link)
	if err != nil {
		t.Fail()
	}
	if link.URL != url {
		t.Fail()
	}
	if link.Slug != nil {
		t.Fail()
	}
}

func TestLinkMarshal(t *testing.T) {
	url := "http://example.com"
	slug := "xyz"
	link := Link{url, &slug}
	bytes, err := json.Marshal(link)
	if err != nil {
		t.Fail()
	}
	str := fmt.Sprintf(`{"url":"%s","slug":"%s"}`, url, slug)
	if str != string(bytes) {
		t.Fail()
	}
}
