package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/oolong-sh/sync"
)

var cfg OolongConfig

type OolongSyncConfig sync.SyncConfig

type OolongConfig struct {
	NotesDirPaths     []string `toml:"note_directories"`
	AllowedExtensions []string `toml:"allowed_extensions"`
	IgnoreDirectories []string `toml:"ignored_directories"`
	OpenCommand       []string `toml:"open_command"`

	// TODO: move these things to a "linker" config section
	NGramRange []int    `toml:"ngram_range"`
	StopWords  []string `toml:"stop_words"`

	PluginsConfig OolongPluginConfig `toml:"plugins"`
	GraphConfig   OolongGraphConfig  `toml:"graph"`
	SyncConfig    OolongSyncConfig   `toml:"sync"`
}

type OolongPluginConfig struct {
	PluginPaths []string `toml:"plugin_paths"`
}

type OolongGraphConfig struct {
	MinNodeWeight float64 `toml:"min_node_weight"`
	MaxNodeWeight float64 `toml:"max_node_weight"`
	MinLinkWeight float64 `toml:"min_link_weight"`
	DefaultMode   string  `toml:"default_mode"`
}

type OolongEditorConfig struct {
	// TODO: web editor related config (themes?)
}

func Config() *OolongConfig { return &cfg }

func NotesDirPaths() []string             { return cfg.NotesDirPaths }
func OpenCommand() []string               { return cfg.OpenCommand }
func NGramRange() []int                   { return cfg.NGramRange }
func AllowedExtensions() []string         { return cfg.AllowedExtensions }
func PluginPaths() []string               { return cfg.PluginsConfig.PluginPaths }
func IgnoredDirectories() []string        { return cfg.IgnoreDirectories }
func StopWords() []string                 { return cfg.StopWords }
func WeightThresholds() OolongGraphConfig { return cfg.GraphConfig }
func GraphMode() string                   { return cfg.GraphConfig.DefaultMode }
func SyncConfig() OolongSyncConfig        { return cfg.SyncConfig }

func Setup() error {
	// CHANGE: for hot-reloading, only one config path location should be supported (~/.config/oolong/oolong.toml)
	configPath, err := findConfigPath()
	if err != nil {
		panic(err)
	}

	readConfig(configPath)

	go initWatcher(configPath)

	return nil
}

func readConfig(configPath string) {
	fmt.Println("Parsing configuration:", configPath)

	contents, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(contents, &cfg)
	if err != nil {
		panic(err)
	}

	for i, dir := range cfg.NotesDirPaths {
		d, err := expand(dir)
		if err != nil {
			panic(err)
		}
		cfg.NotesDirPaths[i] = d
	}

	// TODO: set default values for thresholds if not set
}

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
