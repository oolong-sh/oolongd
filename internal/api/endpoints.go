package api

import (
	"log"
	"net/http"
)

// DOC:
func SpawnServer() {
	mux := http.NewServeMux()

	// TODO: add some sort of JWT system for better security

	// graph endpoints
	mux.HandleFunc("GET /graph", handleGetGraph)

	// config endpoints
	mux.HandleFunc("GET /config", handleGetConfig)
	mux.HandleFunc("GET /config/graph", handleGetGraphConfig)
	mux.HandleFunc("GET /config/note-dirs", handleGetNoteDirsConfig)

	// note endpoints
	mux.HandleFunc("GET /notes", handleGetNotes)
	mux.HandleFunc("GET /note", handleGetNote)
	mux.HandleFunc("POST /note", handleCreateNote)
	mux.HandleFunc("PUT /note", handleUpdateNote)
	mux.HandleFunc("DELETE /note", handleDeleteNote)

	// start server
	log.Println("Starting server on :11975...")
	if err := http.ListenAndServe(":11975", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
