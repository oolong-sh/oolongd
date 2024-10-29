package linking

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/oolong-sh/oolong/internal/linking/lexer"
)

// TODO: multi-document read
// TODO: multi-document ngram merge (maps.copy(dest, src))

var allowedExts = []string{".md", ".mdx", ".tex", ".typ", ".txt"}

// DOC:
func ReadNotesDir(notesDirPaths ...string) ([]*Document, error) {
	documents := []*Document{}

	for _, notesDirPath := range notesDirPaths {
		// expand home dir shorthand
		if strings.HasPrefix(notesDirPath, "~/") || notesDirPath == "~" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			if notesDirPath == "~" {
				notesDirPath = home
			}
			notesDirPath = filepath.Join(home, notesDirPath[2:])
		}

		// extract all note file paths from notes directory
		notePaths := []string{}
		if err := filepath.WalkDir(notesDirPath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			// REFACTOR: probably change this to a blacklist
			// disallow binaries and allow users to blacklist filetypes
			if slices.Contains(allowedExts, filepath.Ext(path)) {
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

	// write out tokens
	// TEST: for debugging, remove later
	b := []byte{}
	for _, d := range documents {
		for _, t := range d.tokens {
			if t.Value == lexer.BreakToken {
				continue
			}
			b = append(b, []byte(t.Lemma+", "+t.Value+"\n")...)
		}
	}
	err := os.WriteFile("./tokens.txt", b, 0666)
	if err != nil {
		panic(err)
	}

	// TEST: for debugging, remove later
	b = []byte{}
	for _, d := range documents {
		for _, ng := range d.ngrams {
			b = append(b, []byte(ng.ngram+"\n")...)
		}
	}
	err = os.WriteFile("./ngrams.txt", b, 0666)
	if err != nil {
		panic(err)
	}

	return documents, nil
}
