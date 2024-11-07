package notes

import (
	"encoding/json"
	"os"

	"github.com/oolong-sh/oolong/internal/documents"
)

var notesFile = "./oolong-notes.json"

type note struct {
	Path    string             `json:"path"`
	Weights map[string]float64 `json:"weights"`
}

// DOC:
func SerializeDocuments(documents map[string]*documents.Document) error {
	notes := documentsToNotes(documents)

	err := serializeNotes(notes)
	if err != nil {
		return err
	}

	return nil
}

func serializeNotes(notes []note) error {
	b, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	err = os.WriteFile(notesFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: parameterize filtering threshold (maybe as a percentage?)
func documentsToNotes(documents map[string]*documents.Document) []note {
	notes := []note{}
	threshold := 8.0

	for k, v := range documents {
		weights := map[string]float64{}
		for k, v := range v.Weights {
			if v > threshold {
				weights[k] = v
			}
		}

		notes = append(notes, note{
			Path:    k,
			Weights: weights,
		})
	}

	return notes
}
