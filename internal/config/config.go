package config

import (
	"os"
	"path/filepath"
)

type OolongConfig struct {
	NotesDirPaths     []string
	NGramRange        []int
	PluginPaths       []string
	AllowedExtensions []string
}

var config OolongConfig

func Config() OolongConfig { return config }

func NotesDirPaths() []string     { return config.NotesDirPaths }
func NGramRange() []int           { return config.NGramRange }
func PluginPaths() []string       { return config.PluginPaths }
func AllowedExtensions() []string { return config.AllowedExtensions }

// TODO: figure out if this should return a mutable pointer or not
// (probably not, use hot reloading based on file modifications)

// CHANGE: possibly use an init method to set this up
// TODO: this function should reload config if it changes
// - watch all entries in .config dir
func Setup(configDir string) (OolongConfig, error) {
	// TODO: Read plugins
	config = OolongConfig{
		NGramRange:        []int{2, 3, 4},
		PluginPaths:       []string{"./scripts/daily_note.lua", "./scripts/event_plugin.lua"},
		NotesDirPaths:     []string{"/home/patrick/notes"},
		AllowedExtensions: []string{".md", ".mdx", ".tex", ".typ", ".txt"},
	}
	// TODO: read config information from lua
	return config, nil
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
