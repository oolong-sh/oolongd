package ngrams

import (
	"slices"
	"strings"

	"github.com/oolong-sh/oolongd/internal/linking/lexer"
)

// add ngram to map, update if it already appears
func addNGram(k string, n int, ngmap map[string]*NGram, i int, tokens []lexer.Lexeme, path string) {
	if ngram, ok := ngmap[k]; ok {
		ngram.globalCount++

		doc := ngram.documents[path]
		doc.DocumentCount++
		// doc.DocumentLocations = append(doc.DocumentLocations, location{row: tokens[i].Row, col: tokens[i].Col})

		// update ngram zone if current is considered more valuable
		if tokens[i].Zone < ngram.zone {
			ngram.zone = tokens[i].Zone
		}
	} else {
		documents := make(map[string]*NGramInfo)

		// create document info struct for ngram
		documents[path] = &NGramInfo{
			DocumentCount:  1,
			DocumentWeight: 0,
			// DocumentLocations: []location{{row: tokens[i].Row, col: tokens[i].Col}},
		}

		// create ngram
		ngmap[k] = &NGram{
			keyword:     k,
			n:           n,
			globalCount: 1,
			documents:   documents,
			zone:        tokens[i].Zone,
		}
	}
}

// Join lexeme items together to form an ngram of length len(nTokens)
// Will return an empty string if stopwords, break tokens, or zone changes occur
func joinNElements(nTokens []lexer.Lexeme) string {
	var parts []string

	// check for outer stop words -> skip ngram
	if slices.Contains(stopWords, strings.ToLower(nTokens[0].Lemma)) ||
		slices.Contains(stopWords, strings.ToLower(nTokens[len(nTokens)-1].Lemma)) {
		return ""
	}

	zone := nTokens[0].Zone

	for i, t := range nTokens {
		// return early if token type matches the break system
		if t.LexType == lexer.Break {
			return ""
		}

		// TODO: add handling of remaining lexeme types (dates)
		if t.LexType == lexer.URI || t.LexType == lexer.Symbol ||
			(t.LexType == lexer.Number && (i == 0 || i == len(nTokens)-1)) {
			return ""
		}

		// return early if tokens are spread across multiple zones
		// TODO: allow zone changes from bold/italic -> default (but not h1 -> default)
		if t.Zone != zone {
			return ""
		}

		// Either value or lemmatized value can be used here, the lemma is likely the better choice
		// parts = append(parts, strings.ToLower(t.Value))
		parts = append(parts, strings.ToLower(t.Lemma))
	}

	return strings.Join(parts, " ")
}
