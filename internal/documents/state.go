package documents

import (
	"fmt"
	"log"
	"os"

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

	//
	// TEST: remove later
	//
	state := State()
	b := append([]byte{}, []byte("ngram,weight,count,ndocs\n")...)
	mng := ngrams.FilterMeaningfulNGrams(state.NGrams, 2, int(float64(len(state.Documents))/1.5), 4.0)
	for _, s := range mng {
		b = append(b, []byte(fmt.Sprintf("%s,%f,%d,%d\n", s, state.NGrams[s].Weight(), state.NGrams[s].Count(), len(state.NGrams[s].Documents())))...)
	}
	if err := os.WriteFile("./meaningful-ngrams.csv", b, 0666); err != nil {
		panic(err)
	}
	//
	// TEST: remove later
	//

	// TODO: other things? (file writes?)

	return nil
}
