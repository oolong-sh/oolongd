package document

import (
	"io"
	"os"
	"strings"
	"unicode"
)

// DOC:
type Lexeme struct {
	Lexeme   string // lexical unit (includes symbols)
	Location int    // row location
}

// DOC:
type Document struct {
	Path     string
	Contents []Lexeme
}

// TODO: figure out if this is actually how I want to handle this
func (d *Document) Text() []Lexeme {
	return d.Contents
}

// DOC:
func ReadDocument(documentPath string) (*Document, error) {
	f, err := os.Open(documentPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := readDocument(f)
	if err != nil {
		return nil, err
	}

	d.Path = documentPath

	return d, nil
}

func readDocument(r io.Reader) (*Document, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	out := &Document{
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
