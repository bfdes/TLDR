package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RedirectHandler expands a shortened link by slug
func RedirectHandler(service LinkService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			msg := http.StatusText(http.StatusMethodNotAllowed)
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		slug := r.URL.Path[1:]
		id, err := Decode(slug)
		if err != nil {
			msg := http.StatusText(http.StatusBadRequest)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		link, found := service.Get(id)
		if found {
			http.Redirect(w, r, link.URL, http.StatusPermanentRedirect)
			return
		}
		http.NotFound(w, r)
	})
}

// CreateLinkHandler generates a new shortened link for its url payload
func CreateLinkHandler(service LinkService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			msg := http.StatusText(http.StatusMethodNotAllowed)
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		if r.Body == nil || r.Body == http.NoBody {
			msg := http.StatusText(http.StatusBadRequest)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		payload := Link{}
		body, err := ioutil.ReadAll(r.Body) // DOS attack vector
		defer r.Body.Close()
		if err == nil {
			err = json.Unmarshal(body, &payload)
		}
		if err != nil {
			msg := http.StatusText(http.StatusBadRequest)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		link, err := service.Create(payload.URL)
		if err != nil {
			msg := http.StatusText(http.StatusInternalServerError)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		res, _ := json.Marshal(link)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	})
}
