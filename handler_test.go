package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type linkServiceStub struct {
	get    func(id int) (Link, bool)
	create func(url string) (Link, error)
}

func (stub linkServiceStub) Get(id int) (link Link, found bool) {
	return stub.get(id)
}

func (stub linkServiceStub) Create(url string) (Link, error) {
	return stub.create(url)
}

func TestRedirect(t *testing.T) {
	url := "http://example.com"
	service := linkServiceStub{
		get: func(id int) (Link, bool) {
			return Link{url, nil}, true
		},
	}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/xyz", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusPermanentRedirect {
		msg := "Unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, http.StatusPermanentRedirect, res.StatusCode)
	}
	if loc := res.Header.Get("Location"); loc != url {
		msg := "Unexpected location: wanted %d, got %d instead"
		t.Errorf(msg, url, loc)
	}
}

func TestRedirectMalformedSlug(t *testing.T) {
	service := linkServiceStub{}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x!z", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	expected := http.StatusBadRequest
	if actual := res.StatusCode; actual != expected {
		msg := "Unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}

func TestRedirectMissingLink(t *testing.T) {
	service := linkServiceStub{
		get: func(id int) (Link, bool) {
			return Link{"", nil}, false
		},
	}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/xyz", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	expected := http.StatusNotFound
	if actual := res.StatusCode; actual != expected {
		msg := "Unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}

func TestCreateLink(t *testing.T) {
	url := "http://example.com"
	slug := "xyz"
	service := linkServiceStub{
		create: func(url string) (Link, error) {
			return Link{url, &slug}, nil
		},
	}
	handler := CreateLinkHandler(service)
	recorder := httptest.NewRecorder()
	payload, _ := json.Marshal(Link{url, nil})
	buffer := bytes.NewBuffer(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/links", buffer)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusCreated {
		msg := "Unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, http.StatusCreated, res.StatusCode)
	}
	link := Link{}
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &link)
	if *link.Slug != slug {
		t.Fail()
	}
}

func TestCreateMalformedPayload(t *testing.T) {
	service := linkServiceStub{}
	handler := CreateLinkHandler(service)
	recorder := httptest.NewRecorder()
	buffer := bytes.NewBufferString(`{"url": "http://example.com"`)
	req, _ := http.NewRequest(http.MethodPost, "/api/links", buffer)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	expected := http.StatusBadRequest
	if actual := res.StatusCode; actual != expected {
		msg := "Unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}

func TestCreateServerError(t *testing.T) {
	service := linkServiceStub{
		create: func(url string) (Link, error) {
			return Link{}, errors.New("db error")
		},
	}
	handler := CreateLinkHandler(service)
	recorder := httptest.NewRecorder()
	buffer := bytes.NewBufferString(`{"url": "http://example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/links", buffer)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	expected := http.StatusInternalServerError
	if actual := res.StatusCode; actual != expected {
		msg := "Unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}
