package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/oolong-sh/sync"
)

var cfg OolongConfig

type OolongSyncConfig sync.SyncConfig

type OolongConfig struct {
	NotesDirPaths     []string `toml:"note_directories"`
	NGramRange        []int    `toml:"ngram_range"`
	AllowedExtensions []string `toml:"allowed_extensions"`
	IgnoreDirectories []string `toml:"ignored_directories"`
	StopWords         []string `toml:"stop_words"`

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
}

func Config() *OolongConfig { return &cfg }

func NotesDirPaths() []string             { return cfg.NotesDirPaths }
func NGramRange() []int                   { return cfg.NGramRange }
func AllowedExtensions() []string         { return cfg.AllowedExtensions }
func PluginPaths() []string               { return cfg.PluginsConfig.PluginPaths }
func IgnoredDirectories() []string        { return cfg.IgnoreDirectories }
func StopWords() []string                 { return cfg.StopWords }
func WeightThresholds() OolongGraphConfig { return cfg.GraphConfig }
func SyncConfig() OolongSyncConfig        { return cfg.SyncConfig }

// TODO: file watcher for config file, reload on change

func Setup() error {
	configPath, err := findConfigPath()
	if err != nil {
		panic(err)
	}

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

	return nil
}
