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
			DocumentWeight:    0, // TODO:
			DocumentLocations: []location{{row: tokens[i].Row, col: tokens[i].Col}},
			DocumentTF:        0, // TODO:
		}
		ngmap[k] = &NGram{
			keyword:      k,
			n:            n,
			globalCount:  1,
			globalWeight: 0, // TODO:
			documents:    documents,
		}
	}
}

// DOC:
func joinNElements(nTokens []lexer.Lexeme) string {
	var parts []string

	// TODO: add handling of different lexeme types (i.e. disallow links)

	// check for outer stop words -> skip ngram
	if slices.Contains(stopWords, nTokens[0].Lemma) ||
		slices.Contains(stopWords, nTokens[len(nTokens)-1].Lemma) {
		return ""
	}

	for _, t := range nTokens {
		// return early if token type matches the break system
		if t.LexType == lexer.Break {
			return ""
		}
		// NOTE: choose between Value and Lemma here
		// parts = append(parts, strings.ToLower(t.Value))
		parts = append(parts, strings.ToLower(t.Lemma))
	}

	return strings.Join(parts, " ")
}
