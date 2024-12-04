package graph

import (
	"encoding/json"
	"fmt"
	"math"
	"path/filepath"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/keywords"
	"github.com/oolong-sh/oolong/internal/notes"
)

type NodeJSON struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Group string  `json:"group"`
	Val   float64 `json:"val"`
}

type LinkJSON struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"strength"`
	Color  string  `json:"color"`
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

const NOTE_NODE_VAL = 1

func SerializeGraph(keywordMap map[string]keywords.Keyword, notes []notes.Note) ([]byte, error) {
	thresholds := config.WeightThresholds()
	minThresh := thresholds.MinNodeWeight
	upperBound := thresholds.MaxNodeWeight

	nodes := []NodeJSON{}
	links := []LinkJSON{}

	for _, keyword := range keywordMap {
		// Only add nodes above the minimum threshold
		if keyword.Weight >= minThresh {
			clampedWeight := clamp(keyword.Weight, minThresh, upperBound)
			nodes = append(nodes, NodeJSON{
				ID:    keyword.Keyword,
				Name:  keyword.Keyword,
				Group: "keyword",
				Val:   clampedWeight,
			})
		}
	}

	for _, note := range notes {
		// Add Note node with a fixed value
		noteID := note.Path
		noteName := filepath.Base(note.Path)
		nodes = append(nodes, NodeJSON{
			ID:    noteID,
			Name:  noteName,
			Group: "note",
			Val:   NOTE_NODE_VAL,
		})

		// Link notes to keywords
		for keywordID, wgt := range note.Weights {
			keyword, exists := keywordMap[keywordID]
			if exists && keyword.Weight >= minThresh {
				links = append(links, LinkJSON{
					Source: noteID,
					Target: keyword.Keyword,
					Value:  wgt,
					Color:  weightToColor(wgt),
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

	// return graph, nil
	return jsonData, nil
}

func weightToColor(v float64) string {
	c := int(math.Round((1 + 204*math.Min(v, 1)) * 1.5))
	return fmt.Sprintf("rgb(%d, %d, %d)", c, c, c)
}
