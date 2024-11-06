package documents

import (
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"sync"

	"github.com/oolong-sh/oolong/internal/config"
)

// DOC: meant to be called with watcher
// assumes paths should not be ignored (should be safe assumption due to watcher ignores)
func ReadDocuments(paths ...string) error {
	// read all input files, update state with documents
	docs := readHandler(paths...)

	// merge ngram maps and calculate weights
	err := updateState(docs)
	if err != nil {
		return err
	}

	// TODO: all weights change, but may not need to be recalculated every time

	return nil
}

// Read, lex, and extract NGrams for all documents in notes directories specified in config file
func ReadNotesDirs() error {
	docs := []*Document{}
	for _, dir := range config.NotesDirPaths() {
		// extract all note file paths from notes directory
		paths := []string{}
		// TODO: add oolong ignore system to blacklist certain subdirs/files
		if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				if slices.Contains(config.IgnoredDirectories(), filepath.Base(path)) {
					return filepath.SkipDir
				}
				return nil
			}

			if slices.Contains(config.AllowedExtensions(), filepath.Ext(path)) {
				paths = append(paths, path)
			}

			return nil
		}); err != nil {
			return err
		}

		// read all documents and append results
		docs = append(docs, readHandler(paths...)...)
	}

	// merge maps and calculate weights
	err := updateState(docs)
	if err != nil {
		return err
	}

	return nil
}

// DOC:
func readHandler(paths ...string) []*Document {
	docs := make([]*Document, len(paths))
	var wg sync.WaitGroup

	// perform a parallel read of found notes files
	wg.Add(len(paths))
	for i, p := range paths {
		go func(i int, notePath string) {
			doc, err := readDocumentByFile(notePath)
			if err != nil {
				log.Printf("Failed to read file: '%s' %v", notePath, err)
				return
			}
			// TODO: this could be changed to use channels
			docs[i] = doc
			wg.Done()
		}(i, p)
	}
	wg.Wait()

	// append results to output array
	return docs
}
