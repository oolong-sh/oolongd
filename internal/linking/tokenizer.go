package linking

import (
	"strings"
	"unicode"
)

// DOC:
type lexeme struct {
	lexeme   string // lexical unit (includes symbols)
	location int    // row location (useful for potential LSP implementation)
}

// DOC:
// TODO: add multiple stages (probably as a param -> int level)
// - (i.e. don't remove capitalization or header tags the first time)
func tokenize(content string) []lexeme {
	out := []lexeme{}
	var sb strings.Builder
	row := 0

	for _, c := range content {
		if unicode.IsSpace(c) {
			if sb.Len() > 0 {
				out = append(out, lexeme{
					lexeme:   sb.String(),
					location: row,
				})
				sb.Reset()
			}

			// FIX: carriage returns need to be handled to avoid incorrect row counts
			if c == '\n' {
				row++
			}
		} else {
			// base case where we want to keep the character
			sb.WriteRune(c)
		}
	}

	// handle remaining content in builder after loop exits
	if sb.Len() > 0 {
		out = append(out, lexeme{
			lexeme:   sb.String(),
			location: row,
		})
	}

	return out
}
