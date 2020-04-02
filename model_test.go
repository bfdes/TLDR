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
	if link.Fragment != nil {
		t.Fail()
	}
}

func TestLinkMarshal(t *testing.T) {
	url := "http://example.com"
	fragment := "xyz"
	link := Link{url, &fragment}
	bytes, err := json.Marshal(link)
	if err != nil {
		t.Fail()
	}
	str := fmt.Sprintf(`{"url":"%s","fragment":"%s"}`, url, fragment)
	if str != string(bytes) {
		t.Fail()
	}
}
