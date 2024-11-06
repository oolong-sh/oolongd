package main

import (
	"flag"
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/daemon"
	"github.com/oolong-sh/oolong/internal/documents"
)

var daemonFlag = flag.Bool("no-daemon", false, "Run Oolong in no-daemon mode (not recommended)")

func main() {
	cfg, err := config.Setup("~/.oolong.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.PluginPaths)

	err = documents.ReadNotesDirs()
	if err != nil {
		return
	}

	// go plugins.InitPlugins(&cfg)
	flag.Parse()
	if !*daemonFlag {
		daemon.Run()
	}
}
