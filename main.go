package main

import (
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/linking"
	"github.com/oolong-sh/oolong/pkg/plugin"
)

func main() {
	cfg, err := config.Setup("./config.lua")
	if err != nil {
		fmt.Println(err)
		return
	}

	d, err := linking.ReadNotesDir("/home/patrick/notes/")
	if err != nil {
		return
	}
	_ = d

	plugin.InitPlugins(&cfg)
}
