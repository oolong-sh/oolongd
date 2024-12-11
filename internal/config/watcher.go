package config

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func initWatcher(configPath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(filepath.Dir(configPath)); err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Watcher event channel returned bad result.")
				return errors.New("Invalid watcher errors channel value.")
			}

			if !strings.Contains(event.Name, configPath) {
				time.Sleep(500)
				continue
			} else if !event.Has(fsnotify.Write) {
				time.Sleep(500)
				continue
			}

			// write event is sent on write start, wait 500ms for write to finish
			time.Sleep(500)
			readConfig(configPath)
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("Invalid watcher errors channel value.")
			}
			log.Println("error:", err)
		}
	}
}
