package state

import (
	"fmt"
	"log"

	"github.com/oolong-sh/oolong/internal/documents"
	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

type OolongState struct {
	Documents map[string]*documents.Document
	NGrams    map[string]*ngrams.NGram
}

type StateManager struct {
	state   OolongState
	updates chan []*documents.Document
	reads   chan chan OolongState
}

var s StateManager

func InitState() {
	s = StateManager{
		state: OolongState{
			Documents: map[string]*documents.Document{},
			NGrams:    map[string]*ngrams.NGram{},
		},
		updates: make(chan []*documents.Document),
		reads:   make(chan chan OolongState),
	}

	go s.run()

	documents.UpdateState = UpdateState
}

func State() OolongState {
	respChan := make(chan OolongState)
	s.reads <- respChan
	log.Println("State fetched.")
	return <-respChan
}

func (s *StateManager) run() {
	for {
		select {
		case docs := <-s.updates:
			s.updateState(docs)
		case resp := <-s.reads:
			resp <- s.state
		}
	}
}

func UpdateState(docs []*documents.Document) error {
	select {
	case s.updates <- docs:
		log.Println("Update request sent")
		return nil
	default:
		log.Println("State update channel is full")
		return fmt.Errorf("state update channel is full")
	}
}

func (s *StateManager) updateState(docs []*documents.Document) {
	log.Println("Updating state and recalculating weights...")

	// Update state documents
	for _, doc := range docs {
		s.state.Documents[doc.Path] = doc
	}

	// Merge resulting n-gram maps
	for _, d := range s.state.Documents {
		ngrams.Merge(s.state.NGrams, d.NGrams)
	}

	// Calculate weights
	ngrams.CalcWeights(s.state.NGrams, len(s.state.Documents))
	log.Println("Done calculating weights.")

	// Update document weights after all weights are calculated
	log.Println("Updating document weights...")
	for ng, ngram := range s.state.NGrams {
		for path, nginfo := range ngram.Documents() {
			s.state.Documents[path].Weights[ng] = nginfo.DocumentWeight
		}
	}
	log.Println("Done updating document weights.")

	//
	// TEST: remove later
	//
	// Generate meaningful n-grams and write to a CSV (for testing/debugging)
	// state := s.state // Snapshot of the current state
	// b := append([]byte{}, []byte("ngram,weight,count,ndocs\n")...)
	// mng := ngrams.FilterMeaningfulNGrams(state.NGrams, 2, int(float64(len(state.Documents))/1.5), 4.0)
	// for _, s := range mng {
	// 	b = append(b, []byte(fmt.Sprintf("%s,%f,%d,%d\n", s, state.NGrams[s].Weight(), state.NGrams[s].Count(), len(state.NGrams[s].Documents())))...)
	// }
	// if err := os.WriteFile("./meaningful-ngrams.csv", b, 0666); err != nil {
	// 	panic(err)
	// }
	// kw := keywords.NGramsToKeywordsMap(s.state.NGrams)
	// notes := notes.DocumentsToNotes(s.state.Documents)
	// dat, err := graph.SerializeGraph(kw, notes)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// if err := os.WriteFile("graph.json", dat, 0644); err != nil {
	// 	panic(err)
	// }
	//
	// TEST: remove later
	//

	log.Println("State update complete.")
}
