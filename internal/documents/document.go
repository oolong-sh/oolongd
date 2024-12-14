package documents

import (
	"io"
	"log"
	"os"

	"github.com/oolong-sh/oolongd/internal/config"
	"github.com/oolong-sh/oolongd/internal/linking/lexer"
	"github.com/oolong-sh/oolongd/internal/linking/ngrams"
)

// Document type stores lexical tokens and NGrams for a single document
type Document struct {
	Path    string
	Weights map[string]float64
	NGrams  map[string]*ngrams.NGram

	tokens []lexer.Lexeme
}

// Read in a single document file, lex, and generate NGrams
// Wraps readDocument for explicit use with files
func readDocumentByFile(documentPath string) (*Document, error) {
	f, err := os.Open(documentPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := readDocument(f, documentPath)
	if err != nil {
		return nil, err
	}
	d.Path = documentPath

	return d, nil
}

// internal reader function that allows usage of io readers for generalized use
func readDocument(r io.Reader, documentPath string) (*Document, error) {
	l := lexer.New()
	log.Printf("Running lexer on %s...\n", documentPath)
	l.Lex(r)

	doc := &Document{
		Path:   documentPath,
		tokens: l.Output,
	}

	// extract ngrams from document
	log.Printf("Generating NGrams for %s...\n", documentPath)
	doc.NGrams = ngrams.Generate(doc.tokens, config.NGramRange(), doc.Path)

	// initialize weights map to avoid nil pointer issues
	doc.Weights = make(map[string]float64, len(doc.NGrams))

	return doc, nil
}
