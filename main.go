package main

import (
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/documents"
)

func main() {
	cfg, err := config.Setup("./config.lua")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.PluginPaths)

	d, err := documents.ReadNotesDir(cfg)
	if err != nil {
		return
	}
	_ = d

	// plugin.InitPlugins(&cfg)
}
