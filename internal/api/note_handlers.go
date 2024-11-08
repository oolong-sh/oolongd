package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/oolong-sh/oolong/internal/state"
)

// 'GET /notes' endpoint handler returns all available note paths
func handleGetNotes(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	data := make([]string, len(state.State().Documents))
	for k := range state.State().Documents {
		data = append(data, k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// DOC:
func handleGetNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter not specified ", http.StatusBadRequest)
		return
	}

	// read file contents
	b, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not read file '"+path+"'", 500)
		return
	}

	// write file contents into response body
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(string(b)); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode file contents as JSON.\n", 500)
	}
}

// DOC:
func handleCreateNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter not specified ", http.StatusBadRequest)
		return
	}

	// check if path before file exists, then check if file exists
	if e, err := exists(path); err != nil {
		log.Println(err)
		http.Error(w, "Error checking path", 500)
		return
	} else if e {
		log.Printf("File %s already exists.\n", path)
		http.Error(w, "Note file already exists", 500)
		return
	}
	// TODO: create directory case?
	// - this will also require adding an additional watcher

	_, err := os.Create(path)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create note file", 500)
		return
	}
}

// DOC:
func handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	// TODO: allow writing notes

	// NOTE: DO NOT CALL STATE UPDATE HERE, LET WATCHER HANDLE IT
	// - this function only needs to write to a file in the watched locations
}

// DOC:
func handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	// TODO: allow writing notes

	// NOTE: this function may need to call the update function due to files no longer existing
	// - check this case in state, this may require substantial logic missing there
}
