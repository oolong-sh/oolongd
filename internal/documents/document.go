package documents

import (
	"fmt"
	"io"
	"os"

	"github.com/oolong-sh/oolong/internal/linking"
	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// TODO: use config
var nGramSizes = []int{2, 3, 4, 5}

// DOC:
type Document struct {
	path   string
	ngwgts map[string]float32 // NOTE: this may need to store more information than just weights

	ngrams map[string]*linking.NGram
	tokens []lexer.Lexeme
}

// Document implementation of Note interface
func (d *Document) NotePath() string             { return d.path }
func (d *Document) Keywords() map[string]float32 { return d.ngwgts }

// DOC:
func ReadDocument(documentPath string) (*Document, error) {
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

// DOC:
func readDocument(r io.Reader, documentPath string) (*Document, error) {
	initStage := 3

	l := lexer.New()
	fmt.Printf("Running lexer on %s...\n", documentPath)
	if err := l.Lex(r, initStage); err != nil {
		return nil, err
	}

	doc := &Document{
		path:   documentPath,
		tokens: l.Output,
	}

	fmt.Printf("Generating NGrams for %s...\n", documentPath)
	doc.ngrams = linking.GenerateNGrams(doc.tokens, nGramSizes, doc.path)

	// for k, v := range out.ngrams {
	// 	fmt.Println("Key:", k, " Value: ", v)
	// }

	// TODO: to avoid re-tokenizing, tokenizer function could take in tokens list
	// after initial stage
	// TODO: multi-stage tokenization and ngram calculation

	doc.setWeightsMap()

	return doc, nil
}

// DOC:
func (d *Document) setWeightsMap() {
	wgts := make(map[string]float32)
	for k, v := range d.ngrams {
		wgts[k] = v.Weight()
	}
	d.ngwgts = wgts
}
