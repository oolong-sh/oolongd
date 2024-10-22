package linking

import (
	"io"
	"os"
)

// DOC:
type Document struct {
	path   string
	ngwgts map[string]float32 // NOTE: this may need to store more information than just weights

	ngrams   []NGram // TODO: may need to be stored as map?
	contents []token
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

func readDocument(r io.Reader) (*Document, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	out := &Document{
		contents: tokenize(string(content)),
	}

	// TODO: multi-stage tokenization and ngram calculation

	return out, nil
}
