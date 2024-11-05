package ngrams

import (
	"fmt"
	"math"
)

// Calculate cosine similarity between two document vectors
func calculateCosineSimilarity(v1, v2 map[string]float64) float64 {
	var dp, mA, mB float64

	for ng, vA := range v1 {
		if vB, ok := v2[ng]; ok {
			dp += vA * vB
		}
		mA += vA * vA
	}

	for _, vB := range v2 {
		mB += vB * vB
	}

	if mA == 0 || mB == 0 {
		return 0
	}

	return dp / (math.Sqrt(mA) * math.Sqrt(mB))
}

// Compute cosine similarity between all document pairs
func CosineSimilarity(ngmap map[string]*NGram) {
	documentVectors := constructDocumentVectors(ngmap)

	for doc1, vec1 := range documentVectors {
		for doc2, vec2 := range documentVectors {
			if doc1 >= doc2 {
				continue
			}
			similarity := calculateCosineSimilarity(vec1, vec2)
			// TODO: do something other than print here? -- (if this actually ends up being used)
			fmt.Printf("%s, %s, %.4f\n", doc1, doc2, similarity)
		}
	}
}

// Construct weighting score vectors
func constructDocumentVectors(ngmap map[string]*NGram) map[string]map[string]float64 {
	documentVectors := make(map[string]map[string]float64)

	for _, ngram := range ngmap {
		for doc, nginfo := range ngram.documents {
			if _, exists := documentVectors[doc]; !exists {
				documentVectors[doc] = make(map[string]float64)
			}
			documentVectors[doc][ngram.keyword] = nginfo.DocumentWeight
		}
	}

	return documentVectors
}
