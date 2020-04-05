package main

// Link domain object
type Link struct {
	URL  string  `json:"url"`
	Slug *string `json:"slug"`
}
