package daemon

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/documents"
)

// Initialize and run file update watcher for notes directories
func runNotesDirsWatcher(dirs ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	dirIgnores := config.IgnoredDirectories()

	for _, dir := range dirs {
		// TODO: add oolong ignore system to blacklist certain subdirs/files
		if err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return nil
			}

			// NOTE: this may not be the exact desired behavior for ignores
			// - this logic also needs to be replicated in the document reader
			if slices.Contains(dirIgnores, filepath.Base(path)) {
				return filepath.SkipDir
			}

			// TEST: this may need to add path as absolute to get correct results
			err = watcher.Add(path)
			if err != nil {
				return err
			}
			log.Println("Added watcher on", path)

			return nil
		}); err != nil {
			return err
		}
	}

	// watcher handler
	// go func() { // running entire function as a goroutine, handler doesn't need to be one
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Watcher event channel returned bad result.")
				return errors.New("Invalid watcher errors channel value.")
			}
			// log.Println("Event:", event)

			if event.Has(fsnotify.Write) {
				log.Println("Modified file:", event.Name)

				// write event is sent on write start, wait 500ms for write to finish
				time.Sleep(500)

				// re-read document
				documents.ReadDocuments(event.Name)
				// TODO: add dedup timer to prevent multi-write calls
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("Invalid watcher errors channel value.")
			}
			log.Println("error:", err)
		}
	}
	// }()
	// <-make(chan struct{})
	// return nil
}
