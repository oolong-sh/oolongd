package ngrams

import (
	"slices"
	"sort"
	"sync"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// DOC:
type NGram struct {
	keyword string
	n       int

	// weight and count across all documents
	globalWeight   float64
	globalCount    int
	inverseDocFreq float64

	// store all documents ngram is present in and info within the document
	documents map[string]*NGramInfo
}

// DOC:
type NGramInfo struct {
	DocumentCount     int
	DocumentWeight    float64
	DocumentLocations []location
	DocumentTermFreq  float64
	DocumentTfIdf     float64
}

type location struct {
	row int
	col int
}

// NGram implements Keyword interface
func (ng *NGram) Weight() float64                  { return ng.globalWeight }
func (ng *NGram) Keyword() string                  { return ng.keyword }
func (ng *NGram) Documents() map[string]*NGramInfo { return ng.documents } // CHANGE: this to return a map of paths to weights?

func (ng *NGram) Count() int   { return ng.globalCount }
func (ng *NGram) IDF() float64 { return ng.inverseDocFreq }

// DOC:
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
				// ngmap[ng].updateWeight() // CHANGE: only calculate weights after maps are merged?
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

	// calculate term frequencies
	tf(ngrams, path)

	return ngrams
}

// TODO: finish this function
func CalcWeights(ngmap map[string]*NGram, N int) {
	idf(ngmap, N)
	tfidf(ngmap)
	// CHANGE: probably take n and word length into account
	for _, ng := range ngmap {
		ng.updateWeight()
	}
}

// TODO: maybe get rid of this function?
func FilterMeaningfulNGrams(ngmap map[string]*NGram, minDF int, maxDF int, minAvgTFIDF float64) []string {
	var out []string
	for k, ng := range ngmap {
		N := len(ng.documents)
		if N >= minDF && N <= maxDF && ng.globalWeight >= minAvgTFIDF {
			out = append(out, k)
		}
	}
	return out
}

// Merge 2 or more string->*NGram maps
func Merge(maps ...map[string]*NGram) {
	for i := 1; i < len(maps); i++ {
		for k, vi := range maps[i] {
			if v0, ok := maps[0][k]; !ok {
				maps[0][k] = vi
			} else {
				v0.globalCount += vi.globalCount
				v0.globalWeight = (v0.globalWeight + vi.globalWeight) / 2 // TODO: more advanced weight logic
				for dk, dv := range vi.documents {
					v0.documents[dk] = dv
				}
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

// TODO: decide what metric to use here (count vs weight vs idf)
func OrderByFrequency(m map[string]*NGram, limit int) []struct {
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

// TODO: functions for filtering less frequent ngrams and stop-words
