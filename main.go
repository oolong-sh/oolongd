package main

import (
	"embed"
	"flag"

	"github.com/oolong-sh/oolongd/internal/config"
	"github.com/oolong-sh/oolongd/internal/daemon"
	"github.com/oolong-sh/oolongd/internal/db"
	"github.com/oolong-sh/oolongd/internal/documents"
	"github.com/oolong-sh/oolongd/internal/linking/ngrams"
	"github.com/oolong-sh/oolongd/internal/state"
)

//go:embed static
var staticFiles embed.FS

var daemonFlag = flag.Bool("no-daemon", false, "Run Oolong in no-daemon mode (not recommended)")

func main() {
	// read in config
	if err := config.Setup(); err != nil {
		panic(err)
	}

	// merge stop words from config
	ngrams.MergeStopWords()

	// initialize state
	state.InitState()

	// read notes directories
	if err := documents.ReadNotesDirs(); err != nil {
		panic(err)
	}

	go func() {
		if config.PinningEnabled() {
			if err := db.InitializeDB(); err != nil {
				panic(err)
			}
			defer db.CloseDB()
		}
	}()

	// go plugins.InitPlugins()

	// run daemon if --no-daemon flag is not passed
	flag.Parse()
	if !*daemonFlag {
		daemon.Run(staticFiles)
	}
}
