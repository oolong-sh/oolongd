package state

import (
	"fmt"
	"log"
	"os"

	"github.com/oolong-sh/oolong/internal/documents"
	"github.com/oolong-sh/oolong/internal/linking/ngrams"
	"github.com/oolong-sh/oolong/pkg/keywords"
	"github.com/oolong-sh/oolong/pkg/notes"
)

// application-wide persistent state of documents and ngrams
var state OolongState

type OolongState struct {
	Documents map[string]*documents.Document
	NGrams    map[string]*ngrams.NGram
}

// State getter
func State() OolongState { return state }

// Initialize oolong state variables and inject state updater function into documents
func InitState() {
	// instantiate persistent state
	state = OolongState{
		Documents: map[string]*documents.Document{},
		NGrams:    map[string]*ngrams.NGram{},
	}

	// dependency injection of state updater function
	documents.UpdateState = UpdateState
}

// Update application state information after file reads are performed
func UpdateState(docs []*documents.Document) error {
	log.Println("Updating state and recalculating weights...")

	// update state documents
	for _, doc := range docs {
		state.Documents[doc.Path] = doc
	}

	// merge resulting ngram maps
	for _, d := range state.Documents {
		ngrams.Merge(state.NGrams, d.NGrams)
	}

	// calculate weights
	ngrams.CalcWeights(state.NGrams, len(state.Documents))
	log.Println("Done calculating weights.")

	// update document weights after all weights are calculated
	log.Println("Updating document weights...")
	for ng, ngram := range state.NGrams {
		for path, nginfo := range ngram.Documents() {
			state.Documents[path].Weights[ng] = nginfo.DocumentWeight
		}
	}
	log.Println("Done updating document weights.")

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

	// serialize results for graph usage
	if err := notes.SerializeDocuments(state.Documents); err != nil {
		panic(err)
	}
	if err := keywords.SerializeNGrams(state.NGrams); err != nil {
		panic(err)
	}

	return nil
}
