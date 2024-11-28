package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oolong-sh/oolong/internal/state"
)

type createUpdateRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// 'GET /notes' endpoint handler returns all available note paths
func handleGetNotes(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	// CORS handling
	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	s := state.State()
	data := []string{}
	for k := range s.Documents {
		data = append(data, k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// 'GET /note?path=<path>' endpoint handler gets note contents corresponding to input path
func handleGetNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	// CORS handling
	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

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

// 'POST /note' endpoint handler creates a note file (and any missing directories) corresponding to input path
// Expected request body: { "path": "/path/to/note", "content", "full note contents to write" }
func handleCreateNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	// CORS handling
	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	// parse request body
	var req createUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode request body", 400)
	}
	log.Println("Request body: ", req)

	// check if path before file exists, then check if file exists
	if e, err := exists(req.Path); err != nil {
		log.Println(err)
		http.Error(w, "Error checking path", 500)
		return
	} else if e {
		log.Printf("File %s already exists.\n", req.Path)
		http.Error(w, "Note file already exists", 500)
		return
	}

	// create directories and file
	dir := filepath.Dir(req.Path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println(err)
		http.Error(w, "Failed to create missing directories", 500)
		return
	}
	if err := os.WriteFile(req.Path, []byte(req.Content), 0644); err != nil {
		log.Println(err)
		http.Error(w, "Failed to create file directories", 500)
		return
	}
}

// 'PUT /note' endpoint handler updates note contents corresponding to input path
// It will create files that do not exist, but will not create directories
// Expected request body: { "path": "/path/to/note", "content", "full note contents to write" }
func handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	// CORS handling
	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	// parse request body
	var req createUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode request body", 400)
	}
	log.Println("Request body: ", req)

	// write contents to file
	if err := os.WriteFile(req.Path, []byte(req.Content), 0666); err != nil {
		log.Println(err)
		http.Error(w, "Failed to write to note file", 500)
		return
	}
}

// 'Delete /note?path=/path/to/note' endpoint handler deletess a note file based on query input
func handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)

	// CORS handling
	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path parameter not specified ", http.StatusBadRequest)
		return
	}

	// attempt to remove file
	if err := os.Remove(path); err != nil {
		log.Println("Failed to delete file", err)
		http.Error(w, "Failed to remove file", 500)
		return
	}

	// NOTE: this function may need to call the update function due to files no longer existing
	// - check this case in state, this may require substantial logic missing there
}
