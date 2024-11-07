package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var config OolongConfig

type OolongConfig struct {
	NotesDirPaths     []string `toml:"note_directories"`
	NGramRange        []int    `toml:"ngram_range"`
	AllowedExtensions []string `toml:"allowed_extensions"`
	PluginPaths       []string `toml:"plugin_paths"`
	IgnoreDirectories []string `toml:"ignored_directories"`
	StopWords         []string `toml:"stopwords"`
}

func Config() *OolongConfig { return &config }

func NotesDirPaths() []string      { return config.NotesDirPaths }
func NGramRange() []int            { return config.NGramRange }
func AllowedExtensions() []string  { return config.AllowedExtensions }
func PluginPaths() []string        { return config.PluginPaths }
func IgnoredDirectories() []string { return config.IgnoreDirectories }
func StopWords() []string          { return config.StopWords }

// TODO: file watcher for config file, reload on change

func Setup(configPath string) (OolongConfig, error) {
	var err error
	configPath, err = expand(configPath)
	if err != nil {
		panic(err)
	}

	contents, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}

	for i, dir := range config.NotesDirPaths {
		d, err := expand(dir)
		if err != nil {
			panic(err)
		}
		config.NotesDirPaths[i] = d
	}

	return config, nil
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
