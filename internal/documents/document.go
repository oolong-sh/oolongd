package documents

import (
	"io"
	"log"
	"os"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/linking/lexer"
	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

// Document type stores lexical tokens and NGrams for a single document
type Document struct {
	path   string
	ngwgts map[string]float64

	ngrams map[string]*ngrams.NGram
	tokens []lexer.Lexeme
}

// Document implementation of Note interface
func (d *Document) Path() string                       { return d.path }
func (d *Document) KeywordWeights() map[string]float64 { return d.ngwgts }

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
	d.path = documentPath

	return d, nil
}

// internal reader function that allows usage of io readers for generalized use
func readDocument(r io.Reader, documentPath string) (*Document, error) {
	l := lexer.New()
	log.Printf("Running lexer on %s...\n", documentPath)
	l.Lex(r)

	doc := &Document{
		path:   documentPath,
		tokens: l.Output,
	}

	log.Printf("Generating NGrams for %s...\n", documentPath)
	doc.ngrams = ngrams.Generate(doc.tokens, config.NGramRange(), doc.path)

	// FIX: weight setting must occur after document NGRam maps are merged
	doc.setWeightsMap()

	return doc, nil
}

// Generate map of weights for all NGrams found in the document
func (d *Document) setWeightsMap() {
	wgts := make(map[string]float64)
	for k, v := range d.ngrams {
		wgts[k] = v.Documents()[d.path].DocumentWeight
	}
	d.ngwgts = wgts
}
