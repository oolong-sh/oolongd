package documents

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/linking/lexer"
	"github.com/oolong-sh/oolong/internal/linking/ngrams"
)

// Read, lex, and extract NGrams for all documents in notes directories specified in config file
func ReadNotesDirs() ([]*Document, error) {
	documents := []*Document{}

	for _, notesDirPath := range config.NotesDirPaths() {
		// extract all note file paths from notes directory
		notePaths := []string{}
		if err := filepath.WalkDir(notesDirPath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			if slices.Contains(config.AllowedExtensions(), filepath.Ext(path)) {
				notePaths = append(notePaths, path)
			}

			return nil
		}); err != nil {
			return nil, err
		}

		// perform a parallel read of found notes files
		var wg sync.WaitGroup
		wg.Add(len(notePaths))
		docs := make([]*Document, len(notePaths))

		for i, notePath := range notePaths {
			go func(i int, notePath string) {
				doc, err := ReadDocument(notePath)
				if err != nil {
					fmt.Printf("Failed to read file: '%s' %v", notePath, err)
					return
				}
				docs[i] = doc
				wg.Done()
			}(i, notePath)
		}

		wg.Wait()

		// append results to output array
		documents = append(documents, docs...)
	}

	//
	// TEST: for debugging, remove later
	//
	// write out tokens
	b := []byte{}
	for _, d := range documents {
		for _, t := range d.tokens {
			if t.Value == lexer.BreakToken {
				continue
			}
			b = append(b, []byte(fmt.Sprintf("%s, %s, %d\n", t.Lemma, t.Value, t.Zone))...)
		}
	}
	err := os.WriteFile("./tokens.txt", b, 0666)
	if err != nil {
		panic(err)
	}

	b = []byte{}
	b = append(b, []byte("ngram,weight,count\n")...)
	ngmap := make(map[string]*ngrams.NGram)
	for _, d := range documents {
		ngrams.Merge(ngmap, d.ngrams)
	}
	ngrams.CalcWeights(ngmap, len(documents))
	for _, d := range documents {
		for _, ng := range d.ngrams {
			b = append(b, []byte(fmt.Sprintf("%s, %f, %d\n", ng.Keyword(), ng.Weight(), ng.Count()))...)
		}
	}
	err = os.WriteFile("./ngrams.txt", b, 0666)
	if err != nil {
		panic(err)
	}
	b = []byte{}
	b = append(b, []byte("ngram,weight,count,ndocs\n")...)
	mng := ngrams.FilterMeaningfulNGrams(ngmap, 2, int(float64(len(documents))/1.5), 4.0)
	for _, s := range mng {
		b = append(b, []byte(fmt.Sprintf("%s,%f,%d,%d\n", s, ngmap[s].Weight(), ngmap[s].Count(), len(ngmap[s].Documents())))...)
	}
	err = os.WriteFile("./meaningful-ngrams.csv", b, 0666)
	if err != nil {
		panic(err)
	}
	// ngrams.CosineSimilarity(ngmap)

	// ngcounts := ngrams.Count(ngmap)
	// freq := ngrams.OrderByFrequency(ngcounts, 10)
	freq := ngrams.OrderByFrequency(ngmap)
	b = []byte{}
	for _, v := range freq {
		b = append(b, []byte(fmt.Sprintf("%s %f\n", v.Key, v.Value))...)
	}
	err = os.WriteFile("./ngram-counts.txt", b, 0666)
	if err != nil {
		panic(err)
	}
	//
	// TEST: for debugging, remove later
	//

	return documents, nil
}
