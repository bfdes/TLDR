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
	get    func(slug string) (string, error)
	create func(url string) (Link, error)
}

func (stub linkServiceStub) Get(slug string) (string, error) {
	return stub.get(slug)
}

func (stub linkServiceStub) Create(url string) (Link, error) {
	return stub.create(url)
}

func TestRedirect(t *testing.T) {
	url := "http://example.com"
	service := linkServiceStub{
		get: func(slug string) (string, error) {
			return url, nil
		},
	}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/xyz", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusPermanentRedirect {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Fatalf(msg, http.StatusPermanentRedirect, res.StatusCode)
	}
	if loc := res.Header.Get("Location"); loc != url {
		msg := "unexpected location: wanted %d, got %d instead"
		t.Errorf(msg, url, loc)
	}
}

func TestRedirectWrongMethod(t *testing.T) {
	service := linkServiceStub{}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/xyz", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Fatalf(msg, http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestRedirectMalformedSlug(t *testing.T) {
	service := linkServiceStub{
		get: func(slug string) (string, error) {
			return "", ErrDecodeFailure
		},
	}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x!z", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	expected := http.StatusBadRequest
	if actual := res.StatusCode; actual != expected {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}

func TestRedirectMissingLink(t *testing.T) {
	service := linkServiceStub{
		get: func(slug string) (string, error) {
			return "", ErrNotFound
		},
	}
	handler := RedirectHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/xyz", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	expected := http.StatusNotFound
	if actual := res.StatusCode; actual != expected {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}

func TestCreate(t *testing.T) {
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
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Fatalf(msg, http.StatusCreated, res.StatusCode)
	}
	link := Link{}
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &link)
	if *link.Slug != slug {
		t.Fail()
	}
}

func TestCreateWrongMethod(t *testing.T) {
	service := linkServiceStub{}
	handler := CreateLinkHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/links", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Fatalf(msg, http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestCreateNoBody(t *testing.T) {
	service := linkServiceStub{}
	handler := CreateLinkHandler(service)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/links", nil)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusBadRequest {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Fatalf(msg, http.StatusBadRequest, res.StatusCode)
	}
}

func TestCreateEmptyBody(t *testing.T) {
	service := linkServiceStub{}
	handler := CreateLinkHandler(service)
	recorder := httptest.NewRecorder()
	buffer := bytes.NewBuffer([]byte{})
	req, _ := http.NewRequest(http.MethodPost, "/api/links", buffer)
	handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusBadRequest {
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Fatalf(msg, http.StatusBadRequest, res.StatusCode)
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
		msg := "unexpected status code: wanted %d, got %d instead"
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
		msg := "unexpected status code: wanted %d, got %d instead"
		t.Errorf(msg, expected, actual)
	}
}
