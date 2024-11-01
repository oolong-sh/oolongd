package ngrams

import (
	"fmt"
	"slices"
	"sort"
	"sync"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// TODO: NGram rework
//
// - need to be able to access ngrams independently of documents (node weight)
// - need to be able to get weights on a per-document basis
//
// - needs to be able to be generated from each document, then merged afterwards
//   - document struct should be able to support either a map[string]*NGram or []NGram
//
// - store keyword/topic/phrase
// - store global weight
// - store total count
// - store map document path -> { document_count, document_weight, document_locations }

// DOC:
type NGram struct {
	keyword string
	n       int

	// weight and count across all documents
	globalWeight float32
	globalCount  int

	// store all documents ngram is present in and info within the document
	documents map[string]*NGramInfo
}

// DOC:
type NGramInfo struct {
	DocumentCount     int
	DocumentWeight    float32
	DocumentLocations []location
}

type location struct {
	row int
	col int
}

// NGram implements Keyword interface
func (ng *NGram) Weight() float32                  { return ng.globalWeight }
func (ng *NGram) Keyword() string                  { return ng.keyword }
func (ng *NGram) Documents() map[string]*NGramInfo { return ng.documents } // CHANGE: this to return a map of paths to weights?

// TODO: update token type to store stage?
// TODO: take in interface of options to show stage, document, stage scaling factor
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
			go func(j int, ngmap map[string]*NGram) {
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
				ngmap[ng].updateWeight(1) // CHANGE: only calculate weights after maps are merged?
			}(j, ngmaps[j])
		}
		wg.Wait()
	}

	// merge result maps
	for _, ngmap := range ngmaps {
		for k, v := range ngmap {
			ngrams[k] = v
		}
	}

	return ngrams
}

// DOC:
// TEST:
func Merge(m1, m2 map[string]*NGram) {
	for k, v2 := range m2 {
		if v1, ok := m1[k]; !ok {
			m1[k] = v2
		} else {
			v1.globalCount += v2.globalCount
			v1.globalWeight = (v1.globalWeight + v2.globalWeight) / 2 // TODO: more advanced weight logic
			for dk, dv := range v2.documents {
				v1.documents[dk] = dv
			}
		}
	}
}

// DOC:
func Count(ngrams map[string]*NGram) map[string]int {
	out := make(map[string]int)

	for k, v := range ngrams {
		out[k] = v.globalCount
	}

	return out
}

// TEST:
func OrderByFrequency(counts map[string]int, limit int) []struct {
	Key   string
	Value int
} {
	kvList := make([]struct {
		Key   string
		Value int
	}, 0, len(counts))

	// Populate the slice with key-value pairs from the map
	for k, v := range counts {
		kvList = append(kvList, struct {
			Key   string
			Value int
		}{k, v})
	}

	// Sort kvList by the values in descending order
	sort.Slice(kvList, func(i, j int) bool {
		return kvList[i].Value > kvList[j].Value
	})

	// Print sorted key-value pairs with values meeting the limit condition
	for _, kv := range kvList {
		if kv.Value >= limit {
			fmt.Printf("%s: %d\n", kv.Key, kv.Value)
		}
	}

	return kvList
}

// TODO: functions for filtering less frequent ngrams and stop-words
