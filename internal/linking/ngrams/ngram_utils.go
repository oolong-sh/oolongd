package ngrams

import (
	"fmt"
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
		// TODO: possibly remove this
		if tokens[i].Zone < ngram.zone {
			ngram.zone = tokens[i].Zone
		}
	} else {
		documents := make(map[string]*NGramInfo)
		documents[path] = &NGramInfo{
			DocumentCount:     1,
			DocumentWeight:    0, // REFACTOR: not currently being used -> use bm25 for this?
			DocumentLocations: []location{{row: tokens[i].Row, col: tokens[i].Col}},
		}
		ngmap[k] = &NGram{
			keyword:     k,
			n:           n,
			globalCount: 1,
			documents:   documents,
			zone:        tokens[i].Zone, // NOTE: maybe get rid of this later
		}
	}
}

// DOC:
func joinNElements(nTokens []lexer.Lexeme) string {
	var parts []string

	// TODO: add handling of different lexeme types (i.e. disallow links)

	// check for outer stop words -> skip ngram
	if slices.Contains(stopWords, strings.ToLower(nTokens[0].Lemma)) ||
		slices.Contains(stopWords, strings.ToLower(nTokens[len(nTokens)-1].Lemma)) {
		return ""
	}
	// for _, t := range nTokens {
	// 	if slices.Contains(stopWords, string(t.Lemma)) {
	// 		return ""
	// 	}
	// }

	zone := nTokens[0].Zone

	for _, t := range nTokens {
		// return early if token type matches the break system
		if t.LexType == lexer.Break {
			return ""
		}
		// return early if tokens are spread across multiple zones
		if t.Zone != zone {
			return ""
		}

		// NOTE: choose between Value and Lemma here
		// parts = append(parts, strings.ToLower(t.Value))
		parts = append(parts, strings.ToLower(t.Lemma))
	}

	out := strings.Join(parts, " ")
	if out == "the author" {
		fmt.Println(nTokens)
	}
	return out
}
