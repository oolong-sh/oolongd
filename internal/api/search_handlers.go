package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/oolong-sh/oolongd/internal/keywords"
	"github.com/oolong-sh/oolongd/internal/notes"
	"github.com/oolong-sh/oolongd/internal/state"
)

func handleSearchKeyword(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "Keyword parameter not specified ", http.StatusBadRequest)
		return
	}

	// TODO: optional query params for min document count/weight

	s := state.State()
	match, exist := keywords.SearchByKeyword(keyword, s.NGrams)
	if !exist {
		log.Printf("Search keyword '%s' not found in corpus.\n", keyword)
		http.Error(w, "Keyword not found in corpus.", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func handleSearchNote(w http.ResponseWriter, r *http.Request) {
	// TODO:

	log.Println("Request received:", r.Method, r.URL, r.Host)

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter not specified ", http.StatusBadRequest)
		return
	}

	// TODO: optional query params for min document count/weight

	s := state.State()
	match, exist := notes.SearchByNote(path, s.Documents)
	if !exist {
		log.Printf("Search path '%s' not found.\n", path)
		http.Error(w, "Note with path not found in corpus.", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}
