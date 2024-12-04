package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oolong-sh/oolong/internal/config"
)

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)
	w.Header().Set("Content-Type", "application/json")

	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	json.NewEncoder(w).Encode(config.Config())
}

func handleGetGraphConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)
	w.Header().Set("Content-Type", "application/json")

	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	json.NewEncoder(w).Encode(config.WeightThresholds())
}

func handleGetNoteDirsConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)
	w.Header().Set("Content-Type", "application/json")

	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	json.NewEncoder(w).Encode(config.NotesDirPaths())
}
