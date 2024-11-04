package main

import (
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/documents"
)

func main() {
	cfg, err := config.Setup("~/.oolong.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.PluginPaths)

	d, err := documents.ReadNotesDir()
	if err != nil {
		return
	}
	_ = d

	// plugins.InitPlugins(&cfg)
}
