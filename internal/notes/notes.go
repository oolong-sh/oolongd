package notes

import (
	"encoding/json"
	"math"

	"github.com/oolong-sh/oolong/internal/config"
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

func DocumentsToNotes(documents map[string]*documents.Document) []Note {
	notes := []Note{}
	minThresh := config.WeightThresholds().MinLinkWeight

	for k, v := range documents {
		weightSum := 0.0
		weights := map[string]float64{}

		// set weight values
		for k, v := range v.Weights {
			if v > minThresh {
				weights[k] = v
				weightSum += v * v
			}
		}

		// normalize resulting weights
		normalizeWeights(weights, math.Sqrt(weightSum))

		notes = append(notes, Note{
			Path:    k,
			Weights: weights,
		})
	}

	return notes
}

func normalizeWeights(m map[string]float64, sum float64) {
	for k, v := range m {
		m[k] = v / sum
	}
}
