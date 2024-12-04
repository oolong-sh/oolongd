package api

import (
	"fmt"
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

	// CORS handling
	if err := checkOrigin(w, r); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintln(err), 500)
		return
	}

	// get snapshot of current state
	s := state.State()

	// convert state into serializable format for graph
	kw := keywords.NGramsToKeywordsMap(s.NGrams)
	notes := notes.DocumentsToNotes(s.Documents)

	// serialize graph data
	data, err := graph.SerializeGraph(kw, notes)
	if err != nil {
		http.Error(w, "Error serializing graph data", 500)
	}

	// encode graph data in reponse
	if _, err := w.Write(data); err != nil {
		http.Error(w, "Error encoding graph data", 500)
		return
	}
}
