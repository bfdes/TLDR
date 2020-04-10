package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestGetUrl(t *testing.T) {
	cache := initCache(cacheHost, cachePort)
	db := initDb(dbHost, dbPort, dbUser, dbPassword, dbName)
	defer db.Close()
	service := linkService{cache, db}

	url, err := uuid.NewRandom()
	if err != nil {
		t.Fatal()
	}
	link, err := service.Create(url.String())
	if err != nil {
		t.Fatal()
	}
	if link.URL != url.String() || link.Slug == nil {
		t.Fatal()
	}

	t.Run("CacheMiss", func(t *testing.T) {
		// n.b. We cannot verify cache hit
		url, err := service.Get(*link.Slug)
		if err != nil {
			msg := "unexpected error: %v"
			t.Fatalf(msg, err)
		}
		if link.URL != url {
			msg := "unexpected url: wanted %s, got %s instead"
			t.Errorf(msg, url, link.URL)
		}
	})
	t.Run("MalformedSlug", func(t *testing.T) {
		url, err := service.Get("x!z")
		if err != ErrDecodeFailure && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err == nil {
			t.Errorf("unexpected url %s", url)
		}
	})
	t.Run("MissingLink", func(t *testing.T) {
		url, err := service.Get("xyz")
		if err != ErrNotFound && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err == nil {
			t.Errorf("unexpected url %s", url)
		}
	})
}

func TestCreateDuplicateLinks(t *testing.T) {
	cache := initCache(cacheHost, cachePort)
	db := initDb(dbHost, dbPort, dbUser, dbPassword, dbName)
	defer db.Close()
	service := linkService{cache, db}

	url, err := uuid.NewRandom()
	if err != nil {
		t.Fatal()
	}
	firstLink, err := service.Create(url.String())
	if err != nil {
		t.Fatalf("unexpected exception: %v", err)
	}
	secondLink, err := service.Create(url.String())
	if err != nil {
		t.Fatalf("unepected exception: %v", err)
	}
	if *firstLink.Slug == *secondLink.Slug {
		msg := "duplicate links for url %s have the same slug %s"
		t.Errorf(msg, *firstLink.Slug, *secondLink.Slug)
	}
}
