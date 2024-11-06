package documents

import (
	"log"

	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

// DOC:
var state = Corpus{
	Documents: map[string]*Document{},
	NGrams:    map[string]*ngrams.NGram{},
}

// DOC:
type Corpus struct {
	Documents map[string]*Document
	NGrams    map[string]*ngrams.NGram
}

// DOC:
func State() Corpus { return state }

// DOC:
func updateState(docs []*Document) error {
	log.Println("Updating state and recalculating weights...")

	// update state documents
	for _, doc := range docs {
		state.Documents[doc.path] = doc
	}

	// merge resulting ngram maps
	for _, d := range state.Documents {
		ngrams.Merge(state.NGrams, d.ngrams)
	}

	// calculate weights
	ngrams.CalcWeights(state.NGrams, len(state.Documents))
	log.Println("Done calculating weights.")

	// TODO: other things? (file writes?)

	return nil
}
