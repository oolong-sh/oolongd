package main

import (
	"flag"
	"fmt"

	"github.com/oolong-sh/oolong/internal/config"
	"github.com/oolong-sh/oolong/internal/daemon"
	"github.com/oolong-sh/oolong/internal/documents"
	"github.com/oolong-sh/oolong/internal/state"
)

var daemonFlag = flag.Bool("no-daemon", false, "Run Oolong in no-daemon mode (not recommended)")

func main() {
	// read in config
	cfg, err := config.Setup("~/.config/oolong.toml")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.PluginPaths)

	// initialize state
	state.InitState()

	// read notes directories
	err = documents.ReadNotesDirs()
	if err != nil {
		panic(err)
	}

	// go plugins.InitPlugins(&cfg)

	// run daemon if --no-daemon flag is not passed
	flag.Parse()
	if !*daemonFlag {
		daemon.Run()
	}
}
