package linking

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

// Characters excluded from special char removal stage
// var allowedSpecialChars = []rune{'\”, '"', '-'}
var allowedSpecialChars = []rune{'-'}

// DOC:
type token struct {
	token    string // lexical unit (includes symbols)
	location int    // row location (useful for potential LSP implementation)
}

// DOC:
// TODO: add multiple stages (probably as a param -> int level)
// - (i.e. don't remove capitalization or header tags the first time)
func tokenize(content string, stage int) []token {
	out := []token{}
	var sb strings.Builder
	row := 0
	hyphenFlag := false

	for _, c := range content {
		c = processChar(c, stage)
		if c == 0 {
			hyphenFlag = false
			if sb.Len() > 0 {
				out = append(out, token{
					token:    sb.String(),
					location: row,
				})
				sb.Reset()
			}

			// NOTE: carriage returns may need to be handled to avoid incorrect row counts
			// FIX: row counter not working?
			if c == '\n' {
				row++
			}
		} else {
			if c == '-' {
				hyphenFlag = true
				continue
			}

			if hyphenFlag && sb.Len() > 0 {
				hyphenFlag = false
				sb.WriteRune('-')
			}

			// TODO: special handling for hyphens (only allow with preceding and succeeding non-whitespace chars)
			sb.WriteRune(c)
		}
	}

	// handle remaining content in builder after loop exits
	if sb.Len() > 0 {
		out = append(out, token{
			token:    sb.String(),
			location: row,
		})
	}

	fmt.Println(out)
	return out
}

// DOC:
func processChar(c rune, stage int) rune {
	if unicode.IsSpace(c) {
		return 0
	}

	// stage 0
	if stage == 0 {
		return c
	}

	// stage 1+
	if stage > 0 {
		c = unicode.ToLower(c)
	}

	// stage 2+
	if stage > 1 {
		if unicode.IsLetter(c) || unicode.IsNumber(c) || slices.Contains(allowedSpecialChars, c) {
			c = unicode.ToLower(c)
		} else {
			return 0
		}
	}

	return c
}
