package linking

import "strings"

// DOC:
type NGram struct {
	ngram  string
	weight float32 // weight of NGram across all documents
	count  int     // count across all documents  NOTE: possibly replace this with a map of ngram->int
	n      int

	// store all documents ngram is present in and counts within the document
	// - maps document path to count of ngram in the document
	// NOTE: this could be changed to []int to store locations
	documents map[string]int

	// TODO: store per-weight documents here?
}

// NGram implements Keyword interface
func (d *NGram) Weight() float32 { return d.weight }
func (d *NGram) Keyword() string { return d.ngram }

// TODO: update token type to store document and stage?
// TODO: take in interface of options to show stage, document, stage scaling factor
func GenerateNGrams(tokens []token, nrange []int) map[string]NGram {
	ngrams := make(map[string]NGram)

	for _, n := range nrange {
		if len(tokens) < n {
			continue
		}

		// TODO: probably move this to be outside loop, loop until max in nrange
		for i := 0; i <= len(tokens)-n; i++ {
			ngString := joinNElements(tokens[i : i+n])

			if ngram, ok := ngrams[ngString]; ok {
				// TODO: add to documents map
				ngram.count++
			} else {
				ngrams[ngString] = NGram{
					ngram:     ngString,
					weight:    float32(n), // TODO: probably needs to be handled elsewhere
					count:     1,
					n:         n,
					documents: map[string]int{}, // TODO: this may not be where this should be handled
				}
			}
		}
	}

	return ngrams
}

func (ng *NGram) updateWeight() {
	// TODO: advanced weight calculations
	// Possible naive formula: (count * n) / (scaling_factor * tokenization_stage)
	// - keep count of total ngrams per document?
	//   - could be used to scale by in-document importance, but might weight against big documents
}

func joinNElements(nTokens []token) string {
	out := ""
	for _, t := range nTokens {
		out = strings.Join([]string{out, t.token}, " ")
	}
	return out
}
