package documents

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/oolong-sh/oolong/internal/config"
)

// State updater function depends on injected function from state to avoid circular dependency
var UpdateState func([]*Document) error

// DOC: meant to be called with watcher
// assumes paths should not be ignored (should be safe assumption due to watcher ignores)
func ReadDocuments(paths ...string) error {
	if UpdateState == nil {
		panic("UpdateState was never instantiated.")
	}

	// read all input files, update state with documents
	docs := readHandler(paths...)

	// merge ngram maps and calculate weights
	err := UpdateState(docs)
	if err != nil {
		return err
	}

	return nil
}

// Read, lex, and extract NGrams for all documents in notes directories specified in config file
func ReadNotesDirs() error {
	if UpdateState == nil {
		panic("UpdateState not instantiated.")
	}

	docs := []*Document{}
	for _, dir := range config.NotesDirPaths() {
		if _, err := os.Stat(dir); err != nil {
			log.Printf("Error reading directory '%s': %v\n", dir, err)
			continue
		}

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

	if len(docs) == 0 {
		return errors.New("No documents found in NotesDirPaths.")
	}

	// merge maps and calculate weights
	err := UpdateState(docs)
	if err != nil {
		return err
	}

	return nil
}

// DOC:
func readHandler(paths ...string) []*Document {
	var wg sync.WaitGroup
	docChan := make(chan *Document)

	// launch a goroutine for each file path and read in parallel
	for _, p := range paths {
		wg.Add(1)
		go func(notePath string) {
			defer wg.Done()
			doc, err := readDocumentByFile(notePath)
			if err != nil {
				log.Printf("Failed to read file: '%s' %v", notePath, err)
				return
			}
			// send the document via channel
			docChan <- doc
		}(p)
	}

	// close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(docChan)
	}()

	// collect documents from the channel
	var docs []*Document
	for doc := range docChan {
		docs = append(docs, doc)
	}

	return docs
}
