package main

import (
	"database/sql"
)

// LinkService fetches and persists links
type LinkService interface {
	Get(id int) (link Link, found bool)
	Create(url string) (link Link, err error)
}

type linkService struct {
	db *sql.DB
}

func (service linkService) Get(id int) (link Link, found bool) {
	query := `
		SELECT url FROM links
		WHERE id=$1
	`
	var url string
	found = true
	err := service.db.QueryRow(query, id).Scan(&url)
	if err != nil {
		found = false // No rows returned or serial overflow
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
	fragment, err := Encode(id)
	return Link{url, &fragment}, err
}
