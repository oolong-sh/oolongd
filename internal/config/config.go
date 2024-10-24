package config

import (
	"os"
	"path/filepath"
)

type oolongConfig struct {
	NGramRange []int
}

var config oolongConfig

// TODO: figure out if this should return a mutable pointer or not
// (probably not, use hot reloading based on file modifications)
func Config() *oolongConfig {
	return &config
}

// CHANGE: possibly use an init method to set this up
// TODO: this function should reload config if it changes
// - watch all entries in .config dir
func Setup(configDir string) error {
	config = oolongConfig{
		NGramRange: []int{2, 3, 4},
	}
	// TODO: read config information from lua
	return nil
}

func initWatcher(configDir string) error {
	filepath.WalkDir(configDir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// TODO: possibly use fsnotify for watching for update events

		return err
	})

	return nil
}
