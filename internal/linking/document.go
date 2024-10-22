package linking

import (
	"fmt"
	"io"
	"os"
)

// DOC:
type Document struct {
	path   string
	ngwgts map[string]float32 // NOTE: this may need to store more information than just weights

	// ngrams []NGram // TODO: may need to be stored as map?
	ngrams map[string]NGram
	tokens []token
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

	d, err := readDocument(f)
	if err != nil {
		return nil, err
	}
	d.path = documentPath

	return d, nil
}

// TODO: functions for filtering less frequent ngrams and stop-words

// DOC:
func readDocument(r io.Reader) (*Document, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// first tokenization stage
	out := &Document{
		tokens: tokenize(string(content)),
	}

	// first ngram extract stage
	out.ngrams = GenerateNGrams(out.tokens, []int{2, 3, 4})
	for k, v := range out.ngrams {
		fmt.Println("Key:", k, " Value: ", v)
	}

	// TODO: multi-stage tokenization and ngram calculation

	// TODO: multi-document ngram merge

	// update map used with interface
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
