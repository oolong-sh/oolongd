package linking

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

// TODO: ngram calculation functions

// TODO: this probably needs to take in a map rather than slice of tokens
// - update token type to store document?
func GenerateNGrams(tokens []token, nrange []int) map[string]NGram {
	ngrams := make(map[string]NGram)

	for _, n := range nrange {
		if len(tokens) < n {
			continue
		}

		// TODO: probably move this to be outside loop, loop until max in nrange
		for i := 0; i <= len(tokens)-n; i++ {
			ngString := extractNElements(tokens[i:i+n], n)

			if ngram, ok := ngrams[ngString]; ok {
				// TODO: add to documents map
				ngram.count++
			} else {
				ngrams[ngString] = NGram{
					ngram:     ngString,
					weight:    1, // TODO: probably needs to be handled elsewhere
					count:     1,
					n:         n,
					documents: map[string]int{}, // TODO: this may not be where this should be handled
				}
			}
		}
	}

	return ngrams
}

func extractNElements(nTokens []token, n int) string {
	out := ""
	for i, t := range nTokens {
		out += t.token
		if i < n-1 {
			out += " "
		}
	}
	return out
}
