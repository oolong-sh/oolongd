package graph

import (
	"encoding/json"
	"path/filepath"

	"github.com/oolong-sh/oolong/internal/keywords"
	"github.com/oolong-sh/oolong/internal/notes"
)

type NodeJSON struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Val  float64 `json:"val"`
}

type LinkJSON struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Graph struct {
	Nodes []NodeJSON `json:"nodes"`
	Links []LinkJSON `json:"links"`
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

const NOTE_NODE_VAL = 10

func SerializeGraph(keywordMap map[string]keywords.Keyword, notes []notes.Note, lowerBound, upperBound float64) ([]byte, error) {
	nodes := []NodeJSON{}
	links := []LinkJSON{}

	for _, keyword := range keywordMap {
		// Only add nodes above the minimum threshold
		if keyword.Weight >= lowerBound {
			clampedWeight := clamp(keyword.Weight, lowerBound, upperBound)
			nodes = append(nodes, NodeJSON{
				ID:   keyword.Keyword,
				Name: keyword.Keyword,
				Val:  clampedWeight,
			})
		}
	}

	for _, note := range notes {
		// Add Note node with a fixed value
		noteID := note.Path
		noteName := filepath.Base(note.Path)
		nodes = append(nodes, NodeJSON{
			ID:   noteID,
			Name: noteName,
			Val:  NOTE_NODE_VAL,
		})

		// Link notes to keywords
		for keywordID := range note.Weights {
			keyword, exists := keywordMap[keywordID]
			if exists && keyword.Weight >= lowerBound {
				links = append(links, LinkJSON{
					Source: noteID,
					Target: keyword.Keyword,
				})
			}
		}
	}

	graph := Graph{
		Nodes: nodes,
		Links: links,
	}

	jsonData, err := json.Marshal(graph)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
