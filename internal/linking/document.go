package linking

import (
	"fmt"
	"io"
	"os"
)

var nGramSizes = []int{2, 3, 4, 5}

// DOC:
type Document struct {
	path   string
	ngwgts map[string]float32 // NOTE: this may need to store more information than just weights

	ngrams map[string]*NGram
	tokens []token
}

// Document implementation of Note interface
func (d *Document) NotePath() string             { return d.path }
func (d *Document) Keywords() map[string]float32 { return d.ngwgts }

// DOC:
func ReadDocument(documentPath string) (*Document, error) {
	// TODO: remove ~/ shorthand in documentPath if present
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

// TODO: functions for filtering less frequent ngrams and stop-words

// DOC:
func readDocument(r io.Reader, documentPath string) (*Document, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// first tokenization stage
	out := &Document{
		path:   documentPath,
		tokens: tokenize(string(content), 0),
	}

	// first ngram extract stage
	out.GenerateNGrams(nGramSizes)
	for k, v := range out.ngrams {
		fmt.Println("Key:", k, " Value: ", v)
	}

	// TODO: to avoid re-tokenizing, tokenizer function could take in tokens list
	// after initial stage
	fmt.Println("Stage 2")
	out.tokens = tokenize(string(content), 1)
	out.GenerateNGrams(nGramSizes)
	for k, v := range out.ngrams {
		fmt.Println("Key:", k, " Value: ", v)
	}

	fmt.Println("Stage 3")
	out.tokens = tokenize(string(content), 2)
	out.GenerateNGrams(nGramSizes)
	for k, v := range out.ngrams {
		fmt.Println("Key:", k, " Value: ", v)
	}

	// TODO: multi-stage tokenization and ngram calculation

	// update map for use with graph-facing interface
	out.setWeightsMap()

	return out, nil
}

// DOC:
func (d *Document) setWeightsMap() {
	wgts := make(map[string]float32)
	for k, v := range d.ngrams {
		wgts[k] = v.weight
	}
	d.ngwgts = wgts
}
