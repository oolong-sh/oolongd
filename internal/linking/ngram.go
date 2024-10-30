package linking

import (
	"slices"
	"strings"
	"sync"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// DOC:
type NGram struct {
	keyword string
	// FIX: add another weight field for weight across all documents
	weight float32 // weight of within a single documents  NOTE: also need weight across documents (map string->f32?)
	count  int     // count within documents  NOTE: possibly replace this with a map of ngram->int
	n      int

	// TODO: store all documents ngram is present in and counts within the document
	document  string
	locations [][2]int
}

// NGram implements Keyword interface
func (ng *NGram) Weight() float32 { return ng.weight } // FIX: should return weight across documents
func (ng *NGram) Keyword() string { return ng.keyword }

// TODO: update token type to store stage?
// TODO: take in interface of options to show stage, document, stage scaling factor
func GenerateNGrams(tokens []lexer.Lexeme, nrange []int, path string) map[string]*NGram {
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
				ngString := joinNElements(tokens[i : i+n])
				if ngString == "" {
					return
				}

				// check if ngram is already present in map
				if ngram, ok := ngmap[ngString]; ok {
					ngram.count++
					ngram.locations = append(ngram.locations, [2]int{tokens[i].Row, tokens[i].Col})
				} else {
					ngmap[ngString] = &NGram{
						keyword:   ngString,
						count:     1,
						n:         n,
						document:  path,
						locations: [][2]int{{tokens[i].Row, tokens[i].Col}},
					}
				}

				ngmap[ngString].updateWeight(1)
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

	// TODO:
	return ngrams
}

// DOC:
func joinNElements(nTokens []lexer.Lexeme) string {
	out := ""

	// TODO: add handling of different lexeme types (i.e. disallow links)

	// check for outer stop words -> skip ngram
	if slices.Contains(stopWords, nTokens[0].Value) ||
		slices.Contains(stopWords, nTokens[len(nTokens)-1].Value) {
		return ""
	}

	for _, t := range nTokens {
		// return early if token type matches the break system
		if t.LexType == lexer.Break {
			return ""
		}

		// TODO: make number of stopwords count toward the weight negatively?
		out = strings.Join([]string{out, t.Value}, " ")
	}
	return out
}

// DOC:
func CountNGrams(ngrams map[string]*NGram) map[string]int {
	// TODO:
	out := make(map[string]int)

	for k, _ := range ngrams {
		if _, ok := out[k]; !ok {
			out[k] = 1
		} else {
			out[k]++
		}
	}

	return out
}

// TODO: functions for filtering less frequent ngrams and stop-words

// DOC:
func (ng *NGram) updateWeight(stage int) {
	countWeighting := 0.8 * float32(ng.count)
	nWeighting := 0.3 * float32(ng.n)
	stageWeighting := 0.5 * (float32(stage) + 0.01)

	// TODO: advanced weight calculations
	// Possible naive formula: (count * n) / (scaling_factor * tokenization_stage)
	// - keep count of total ngrams per document?
	//   - could be used to scale by in-document importance, but might weight against big documents
	ng.weight = (countWeighting + nWeighting) / (stageWeighting)
}
