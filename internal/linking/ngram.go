package linking

import (
	"slices"
	"strings"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// DOC:
type NGram struct {
	ngram string
	// FIX: add another weight field for weight across all documents
	weight float32 // weight of within a single documents NOTE: need one across documents
	count  int     // count across all documents  NOTE: possibly replace this with a map of ngram->int
	n      int

	// TODO: store all documents ngram is present in and counts within the document
	document  string
	locations [][2]int
}

// NGram implements Keyword interface
func (ng *NGram) Weight() float32 { return ng.weight } // FIX: should return weight across documents
func (ng *NGram) Keyword() string { return ng.ngram }

// TODO: update token type to store stage?
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
				ngram.locations = append(ngram.locations, d.tokens[i].Location)
			} else {
				ngrams[ngString] = &NGram{
					ngram:     ngString,
					count:     1,
					n:         n,
					document:  d.path,
					locations: [][2]int{d.tokens[i].Location},
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
func joinNElements(nTokens []lexer.Lexeme) string {
	out := ""

	// check for outer stop words -> skip ngram
	if slices.Contains(stopWords, nTokens[0].Value) ||
		slices.Contains(stopWords, nTokens[len(nTokens)-1].Value) {
		return ""
	}

	for _, t := range nTokens {
		// return early if tokens slice contains break sequence
		if t.Value == lexer.BreakToken {
			return ""
		}

		// TODO: handle stop words (but allow in the middle of the word)
		// - make number of stopwords count toward the weight negatively?
		out = strings.Join([]string{out, t.Value}, " ")
	}
	return out
}

// TODO: add smart filtering system for tokens
// - need to be able to filter out noisy tokens
// - could use some sort of ml validation or a dictionary
