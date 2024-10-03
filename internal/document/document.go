package document

import (
	"io"
	"strings"
	"unicode"
)

type Lexeme struct {
	Lexeme   string // lexical unit
	Location int    // row location
}

type Document struct {
	Path     string
	Contents []Lexeme
}

func ReadDocument(r io.Reader, documentPath string) (*Document, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	out := &Document{
		Path:     documentPath,
		Contents: readContents(string(content)),
	}

	return out, nil
}

func readContents(content string) []Lexeme {
	out := []Lexeme{}
	var sb strings.Builder
	row := 0

	for _, c := range content {
		if unicode.IsSpace(c) {
			if sb.Len() > 0 {
				out = append(out, Lexeme{
					Lexeme:   sb.String(),
					Location: row,
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
		out = append(out, Lexeme{
			Lexeme:   sb.String(),
			Location: row,
		})
	}

	return out
}
