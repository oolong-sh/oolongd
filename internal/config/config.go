package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var cfg OolongConfig

type OolongConfig struct {
	NotesDirPaths     []string `toml:"note_directories"`
	NGramRange        []int    `toml:"ngram_range"`
	AllowedExtensions []string `toml:"allowed_extensions"`
	PluginPaths       []string `toml:"plugin_paths"`
	IgnoreDirectories []string `toml:"ignored_directories"`
	StopWords         []string `toml:"stop_words"`
}

func Config() *OolongConfig { return &cfg }

func NotesDirPaths() []string      { return cfg.NotesDirPaths }
func NGramRange() []int            { return cfg.NGramRange }
func AllowedExtensions() []string  { return cfg.AllowedExtensions }
func PluginPaths() []string        { return cfg.PluginPaths }
func IgnoredDirectories() []string { return cfg.IgnoreDirectories }
func StopWords() []string          { return cfg.StopWords }

// TODO: file watcher for config file, reload on change

func Setup(configPath string) error {
	var err error
	configPath, err = expand(configPath)
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

	return nil
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
