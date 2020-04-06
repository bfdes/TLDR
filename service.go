package main

import (
	"database/sql"
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

// LinkService fetches and persists links
type LinkService interface {
	Get(slug string) (string, error)
	Create(url string) (Link, error)
}

type linkService struct {
	cache *memcache.Client
	db    *sql.DB
}

// ErrDecodeFailure is returned by Get when the slug cannot be decoded to an ID
var ErrDecodeFailure = errors.New("cannot decode slug")

// ErrNotFound is returned by Get when a link cannot be found
var ErrNotFound = errors.New("link not found")

func (service linkService) Get(slug string) (string, error) {
	item, err := service.cache.Get(slug)
	if err == nil {
		// Cache hit
		return string(item.Value), nil
	}
	// Cache miss or malformed key
	id, err := Decode(slug)
	if err != nil {
		return "", ErrDecodeFailure
	}
	var url string
	query := `
		SELECT url FROM links
		WHERE id=$1
	`
	err = service.db.QueryRow(query, id).Scan(&url)
	if err != nil {
		// No rows returned or serial overflow
		return "", ErrNotFound
	}
	// Write to cache, but serve request even on error
	item = &memcache.Item{Key: slug, Value: []byte(url)}
	service.cache.Set(item)
	return url, nil
}

func (service linkService) Create(url string) (Link, error) {
	query := `
		INSERT INTO links(url)
		VALUES ($1)
		RETURNING id
	`
	id := 0
	err := service.db.QueryRow(query, url).Scan(&id)
	if err != nil {
		return Link{}, err
	}
	slug, err := Encode(id)
	return Link{url, &slug}, err
}
