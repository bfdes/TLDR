package main

import (
	"database/sql"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

// LinkService fetches and persists links
type LinkService interface {
	Get(id int) (link Link, found bool)
	Create(url string) (link Link, err error)
}

type linkService struct {
	cache *memcache.Client
	db    *sql.DB
}

func (service linkService) Get(id int) (link Link, found bool) {
	var url string
	strID := strconv.Itoa(id)
	item, err := service.cache.Get(strID)

	if err == nil {
		// Cache hit
		found = true
		url = string(item.Value)
	} else {
		// Cache miss
		query := `
			SELECT url FROM links
			WHERE id=$1
		`
		err = service.db.QueryRow(query, id).Scan(&url)
		if err == nil {
			// Write to cache, but serve request even on error
			found = true
			item = &memcache.Item{Key: strID, Value: []byte(url)}
			service.cache.Set(item)
		}
		// No rows returned or serial overflow
	}
	link = Link{url, nil} // Don't bother re-encoding the ID
	return
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
