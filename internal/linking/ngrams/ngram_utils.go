package ngrams

import (
	"slices"
	"strings"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// add ngram to map, update if it already appears
func addNGram(k string, n int, ngmap map[string]*NGram, i int, tokens []lexer.Lexeme, path string) {
	if ngram, ok := ngmap[k]; ok {
		ngram.globalCount++
		doc := ngram.documents[path]
		doc.DocumentCount++
		doc.DocumentLocations = append(doc.DocumentLocations, location{row: tokens[i].Row, col: tokens[i].Col})
	} else {
		documents := make(map[string]*NGramInfo)
		documents[path] = &NGramInfo{
			DocumentCount:     1,
			DocumentLocations: []location{{row: tokens[i].Row, col: tokens[i].Col}},
		}
		ngmap[k] = &NGram{
			keyword:      k,
			n:            n,
			globalCount:  1,
			globalWeight: 0,
			documents:    documents,
		}
	}
}

// DOC:
func joinNElements(nTokens []lexer.Lexeme) string {
	var out string

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
		out = strings.Join([]string{out, strings.ToLower(t.Value)}, " ")
		// out = strings.Join([]string{out, t.Value}, " ")
	}
	// FIX: out has a leading whitespace

	return out
}

// DOC:
func (ng *NGram) updateWeight(stage int) {
	countWeighting := 0.8 * float32(ng.globalCount)
	nWeighting := 0.3 * float32(ng.n)
	stageWeighting := 0.5 * (float32(stage) + 0.01)

	// TODO: advanced weight calculations
	// Possible naive formula: (count * n) / (scaling_factor * tokenization_stage)
	// - keep count of total ngrams per document?
	//   - could be used to scale by in-document importance, but might weight against big documents
	ng.globalWeight = (countWeighting + nWeighting) / (stageWeighting)
}
