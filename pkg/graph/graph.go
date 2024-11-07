package graph

import (
	"encoding/json"
	"path/filepath"

	"github.com/oolong-sh/oolong/pkg/keywords"
	"github.com/oolong-sh/oolong/pkg/notes"
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

const NOTE_NODE_VAL = 50

func SerializeGraph(keywordMap map[string]keywords.Keyword, notes []notes.Note, lowerBound, upperBound float64) ([]byte, error) {
	nodes := []NodeJSON{}
	links := []LinkJSON{}

	for _, keyword := range keywordMap {
		clampedWeight := clamp(keyword.Weight, lowerBound, upperBound)
		nodes = append(nodes, NodeJSON{
			ID:   keyword.Keyword,
			Name: keyword.Keyword,
			Val:  clampedWeight,
		})
	}

	for _, note := range notes {
		// Add Note node
		noteID := note.Path
		noteName := filepath.Base(note.Path) // /home/patrick/notes/home/blogs/bayes.md -> bayes.md
		nodes = append(nodes, NodeJSON{
			ID:   noteID,
			Name: noteName,
			Val:  NOTE_NODE_VAL,
		})

		// Link notes to keywords
		for keywordID := range note.Weights {
			keyword, exists := keywordMap[keywordID]
			if exists {
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
