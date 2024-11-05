package daemon

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Initialize and run file update watcher for notes directories
func runNotesDirsWatcher(dirs ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for _, dir := range dirs {
		// TODO: add oolong ignore system to blacklist certain subdirs/files
		if err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				err = watcher.Add(path)
				if err != nil {
					return err
				}
				log.Println("Added watcher on", path)
			}
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
			log.Println("Event:", event)

			// TODO: call functions on document update
			// - has potential performance impacts, so only rerun what is needed (i.e. lex single doc + weight updates)
			if event.Has(fsnotify.Write) {
				log.Println("Modified file:", event.Name)
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
