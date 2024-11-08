package api

import (
	"encoding/json"
	"net/http"
)

func handleGetGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: finish implementation
	// data, err := graph.SerializeGraph()
	data := ""

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding graph data", 400)
	}

}
