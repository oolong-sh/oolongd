package main

import (
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/linking"
)

func main() {
	cfg, err := config.Setup("./config.lua")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.PluginPaths)

	d, err := linking.ReadNotesDir(cfg.NotesDirPaths...)
	if err != nil {
		return
	}
	_ = d

	// plugin.InitPlugins(&cfg)
}
