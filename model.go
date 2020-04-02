package main

// Link domain object
type Link struct {
	URL      string  `json:"url"`
	Fragment *string `json:"fragment"`
}
