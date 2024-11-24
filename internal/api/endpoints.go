package api

import (
	"log"
	"net/http"
)

// DOC:
func SpawnServer() {
	mux := http.NewServeMux()

	// graph endpoints
	mux.HandleFunc("GET /graph", handleGetGraph)

	// note endpoints
	mux.HandleFunc("GET /notes", handleGetNotes)
	mux.HandleFunc("GET /note", handleGetNote)
	mux.HandleFunc("POST /note", handleCreateNote)
	mux.HandleFunc("PUT /note", handleUpdateNote)
	mux.HandleFunc("DELETE /note", handleDeleteNote)

	// keyword endpoints?

	// plugin endpoints? (probably not these outside of official ones)

	// start server
	log.Println("Starting server on :11975...")
	if err := http.ListenAndServe(":11975", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
