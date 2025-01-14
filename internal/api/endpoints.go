package api

import (
	"embed"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/oolong-sh/oolongd/internal/config"
)

// spawn the oolong api server
func SpawnServer(staticFiles embed.FS) {
	mux := http.NewServeMux()

	// graph endpoints
	mux.HandleFunc("GET /graph", handleGetGraph)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			data, err := staticFiles.ReadFile("static/index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
			return
		}

		// Serve other static files (e.g., .css, .js)
		filePath := "static" + r.URL.Path
		data, err := staticFiles.ReadFile(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Determine the MIME type based on the file extension
		ext := filepath.Ext(filePath)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	})

	// config endpoints
	mux.HandleFunc("GET /config", handleGetConfig)
	mux.HandleFunc("GET /config/graph", handleGetGraphConfig)
	mux.HandleFunc("GET /config/note-dirs", handleGetNoteDirsConfig)
	mux.HandleFunc("GET /config/default-graph-mode", handleGetGraphView)

	// note endpoints
	mux.HandleFunc("GET /notes", handleGetNotes)
	mux.HandleFunc("GET /note", handleGetNote)
	mux.HandleFunc("POST /note", handleCreateNote)
	mux.HandleFunc("PUT /note", handleUpdateNote)
	mux.HandleFunc("DELETE /note", handleDeleteNote)
	mux.HandleFunc("GET /open/note", handleOpenNote)

	// pinning endpoints
	if config.PinningEnabled() {
		mux.HandleFunc("GET /pins", handleGetPinnedNotes)
		mux.HandleFunc("POST /pins", handleAddPinnedNote)
		mux.HandleFunc("DELETE /pins", handleDeletePinnedNote)
	}

	// search endpoints
	mux.HandleFunc("GET /search/keyword", handleSearchKeyword)
	mux.HandleFunc("GET /search/note", handleSearchNote)

	// start server
	log.Println("Starting server on :11975...")
	if err := http.ListenAndServe(":11975", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
