package daemon

import (
	"github.com/oolong-sh/oolong/internal/api"
	"github.com/oolong-sh/oolong/internal/config"
)

// Launch perpetually running watchers and run application in the background as a daemon
func Run() {
	// run file watcher
	go runNotesDirsWatcher(config.NotesDirPaths()...)

	// run api server
	go api.SpawnServer()

	// run forever
	<-make(chan struct{})
}
