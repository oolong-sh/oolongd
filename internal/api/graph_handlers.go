package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/oolong-sh/oolong/internal/graph"
	"github.com/oolong-sh/oolong/internal/keywords"
	"github.com/oolong-sh/oolong/internal/notes"
	"github.com/oolong-sh/oolong/internal/state"
)

func handleGetGraph(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL, r.Host)
	w.Header().Set("Content-Type", "application/json")

	// get snapshot of current state
	s := state.State()

	// convert state into serializable format for graph
	kw := keywords.NGramsToKeywordsMap(s.NGrams)
	notes := notes.DocumentsToNotes(s.Documents)

	// serialize graph data
	// TODO: pass in thresholds (with request? or with config?)
	data, err := graph.SerializeGraph(kw, notes, 0.1, 80)
	if err != nil {
		http.Error(w, "Error serializing graph data", 500)
	}

	// encode graph data in reponse
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding graph data", 500)
	}
}
