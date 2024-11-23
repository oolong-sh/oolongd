package notes

import (
	"encoding/json"

	"github.com/oolong-sh/oolong/internal/documents"
)

type Note struct {
	Path    string             `json:"path"`
	Weights map[string]float64 `json:"weights"`
}

// DOC:
func SerializeDocuments(documents map[string]*documents.Document) ([]byte, error) {
	notes := DocumentsToNotes(documents)

	b, err := json.Marshal(notes)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// TODO: parameterize filtering threshold (maybe as a percentage?)
func DocumentsToNotes(documents map[string]*documents.Document) []Note {
	notes := []Note{}
	threshold := 2.0

	for k, v := range documents {
		weights := map[string]float64{}
		for k, v := range v.Weights {
			if v > threshold {
				weights[k] = v
			}
		}

		notes = append(notes, Note{
			Path:    k,
			Weights: weights,
		})
	}

	return notes
}
