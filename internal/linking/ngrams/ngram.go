package ngrams

import (
	"slices"
	"sort"
	"sync"

	"github.com/oolong-sh/oolongd/internal/linking/lexer"
)

// NGram type used throughout linking package
type NGram struct {
	keyword string
	n       int

	// weight and count across all documents
	globalWeight float64
	globalCount  int
	idf          float64
	zone         lexer.Zone

	// store all documents ngram is present in and info within the document
	documents map[string]*NGramInfo
}

// Information about NGram occurences in a single document
type NGramInfo struct {
	DocumentCount  int
	DocumentWeight float64
	// DocumentLocations []location
	// DocumentTF        float64
	// DocumentTfIdf     float64
	// DocumentBM25      float64
}

// location type for occurence of an NGram within a document
type location struct {
	row int
	col int
}

// NGram implements Keyword interface
func (ng *NGram) Weight() float64 { return ng.globalWeight }
func (ng *NGram) Keyword() string { return ng.keyword }

// Non-interface getter methods
func (ng *NGram) Count() int                       { return ng.globalCount }
func (ng *NGram) Documents() map[string]*NGramInfo { return ng.documents }

// func (ng *NGram) IDF() float64                     { return ng.idf }

// Generate NGrams from a slice of document lexemes
func Generate(tokens []lexer.Lexeme, nrange []int, path string) map[string]*NGram {
	ngrams := make(map[string]*NGram)
	slices.Sort(nrange)

	// set up parallelization variables
	var wg sync.WaitGroup
	ngmaps := make([]map[string]*NGram, len(nrange))
	for i := range ngmaps {
		ngmaps[i] = make(map[string]*NGram)
	}

	// iterate over all tokens in document
	for i := 0; i <= len(tokens)-nrange[0]; i++ {
		// iterate over each size of N
		wg.Add(len(nrange))
		for j, n := range nrange {
			// check for ngrams of each size in nrange in parallel
			go func(j int, n int, ngmap map[string]*NGram) {
				defer wg.Done()
				if i+n > len(tokens) {
					return
				}

				// get string representation of ngram string
				ng := joinNElements(tokens[i : i+n])
				if ng == "" {
					return
				}

				// check if ngram is already present in map
				addNGram(ng, n, ngmap, i, tokens, path)
			}(j, n, ngmaps[j])
		}
		wg.Wait()
	}

	// merge result maps
	for _, ngmap := range ngmaps {
		for k, v := range ngmap {
			ngrams[k] = v
		}
	}

	// calculate term frequencies
	// tf(ngrams, path) // using count instead

	return ngrams
}

// Merge 2 or more string->*NGram maps
func Merge(maps ...map[string]*NGram) {
	for i := 1; i < len(maps); i++ {
		for k, vi := range maps[i] {
			if v0, ok := maps[0][k]; !ok {
				// ngram key not found in main map, add it
				maps[0][k] = vi
			} else {
				// ngram key found in map, merge counts and document info
				// weights should be calculated elsewhere after all merges are completed
				v0.globalCount += vi.globalCount
				// lower zones are considered better, take best
				if v0.zone > vi.zone {
					v0.zone = vi.zone
				}
				for dk, dv := range vi.documents {
					v0.documents[dk] = dv
				}
			}
		}
	}
}

// Count occurences of each NGram in a map of NGrams
func Count(ngrams map[string]*NGram) map[string]int {
	out := make(map[string]int)

	for k, v := range ngrams {
		out[k] = v.globalCount
	}

	return out
}

// Create an slice of NGrams keywords, ordered by their weights
func OrderByFrequency(m map[string]*NGram) []struct {
	Key   string
	Value float64
} {
	kvList := make([]struct {
		Key   string
		Value float64
	}, 0, len(m))

	// Populate the slice with key-value pairs from the map
	for k, v := range m {
		kvList = append(kvList, struct {
			Key   string
			Value float64
		}{k, v.globalWeight})
	}

	// Sort kvList by the values in descending order
	sort.Slice(kvList, func(i, j int) bool {
		return kvList[i].Value < kvList[j].Value
	})

	return kvList
}
