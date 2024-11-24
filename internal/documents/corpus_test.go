package documents_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/documents"
	"github.com/oolong-sh/oolong/internal/linking/ngrams"
	"github.com/oolong-sh/oolong/internal/state"
)

var cfg = config.OolongConfig{
	// NotesDirPaths:     []string{"~/notes"},
	NotesDirPaths:     []string{"../../examples/data"},
	NGramRange:        []int{1, 2, 3},
	AllowedExtensions: []string{".md", ".tex", ".typ", ".txt"},
	PluginPaths:       []string{},
	IgnoreDirectories: []string{".templates", ".git"},
}

func init() {
	state.InitState()
}

func TestReadNotesDirs(t *testing.T) {
	conf := config.Config()
	conf.NotesDirPaths = cfg.NotesDirPaths
	conf.NGramRange = cfg.NGramRange
	conf.AllowedExtensions = cfg.AllowedExtensions
	conf.PluginPaths = cfg.PluginPaths
	conf.IgnoreDirectories = cfg.IgnoreDirectories

	// TODO: actual tests with an example data directory
	fmt.Println("reading?? -- gets lock")
	if err := documents.ReadNotesDirs(); err != nil {
		t.Fatalf("Failed to read notes directories: %v\n", err)
	}
	fmt.Println("finished reading -- getting read lock")
	s := state.State()

	b := append([]byte{}, []byte("ngram,weight,count\n")...)
	for _, d := range s.Documents {
		for _, ng := range d.NGrams {
			b = append(b, []byte(fmt.Sprintf("%s, %f, %d\n", ng.Keyword(), ng.Weight(), ng.Count()))...)
		}
	}
	if err := os.WriteFile("./ngrams.txt", b, 0666); err != nil {
		t.Fatalf("Failed to write ngrams file: %v\n", err)
	}

	b = append([]byte{}, []byte("ngram,weight,count,ndocs\n")...)
	mng := ngrams.FilterMeaningfulNGrams(s.NGrams, 2, int(float64(len(s.Documents))/1.5), 4.0)
	for _, k := range mng {
		b = append(b, []byte(fmt.Sprintf("%s,%f,%d,%d\n", k, s.NGrams[k].Weight(), s.NGrams[k].Count(), len(s.NGrams[k].Documents())))...)
	}
	if err := os.WriteFile("./meaningful-ngrams.csv", b, 0666); err != nil {
		t.Fatalf("Failed to write out meaningful ngrams: %v\n", err)
	}

	// ngrams.CosineSimilarity(ngmap)
}
