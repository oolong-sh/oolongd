package linking

import (
	"slices"
	"strings"
)

// DOC:
type NGram struct {
	ngram  string
	weight float32 // weight of NGram across all documents
	count  int     // count across all documents  NOTE: possibly replace this with a map of ngram->int
	n      int

	// store all documents ngram is present in and counts within the document
	document  string
	locations []int
}

// NGram implements Keyword interface
func (ng *NGram) Weight() float32 { return ng.weight }
func (ng *NGram) Keyword() string { return ng.ngram }

// TODO: update token type to store document and stage?
// TODO: take in interface of options to show stage, document, stage scaling factor
func (d *Document) GenerateNGrams(nrange []int) {
	ngrams := make(map[string]*NGram)

	slices.Sort(nrange)

	// iterate over all tokens in document
	for i := 0; i <= len(d.tokens)-nrange[0]; i++ {
		// iterate over each size of N
		for _, n := range nrange {
			if i+n > len(d.tokens) {
				break
			}

			// get string representation of ngram string
			ngString := joinNElements(d.tokens[i : i+n])
			if ngString == "" {
				continue
			}

			// check if ngram is already present in map
			if ngram, ok := ngrams[ngString]; ok {
				ngram.count++
				ngram.locations = append(ngram.locations, d.tokens[i].location)
			} else {
				ngrams[ngString] = &NGram{
					ngram:     ngString,
					count:     1,
					n:         n,
					document:  d.path,
					locations: []int{d.tokens[i].location},
				}
			}

			ngrams[ngString].updateWeight(1)
		}
	}

	d.ngrams = ngrams
}

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

// DOC:
func joinNElements(nTokens []token) string {
	out := ""

	// check for outer stop words -> skip ngram
	if slices.Contains(stopWords, nTokens[0].token) ||
		slices.Contains(stopWords, nTokens[len(nTokens)-1].token) {
		return out
	}
	// CHANGE: tokenizer needs some way of indicating where special chars were
	// - need to know to use them as stop chars in many cases?
	// - also needs to disallow hyphens if they aren't directly surrounded by other valid chars

	for _, t := range nTokens {
		// TODO: handle stop words (but allow in the middle of the word)
		// - make number of stopwords count toward the weight negatively?
		out = strings.Join([]string{out, t.token}, " ")
	}
	return out
}
