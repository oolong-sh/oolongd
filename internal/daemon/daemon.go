package daemon

import "github.com/oolong-sh/oolong/internal/config"

// Launch perpetually running watchers and run application in the background as a daemon
func Run() {
	go runNotesDirsWatcher(config.NotesDirPaths()...)

	// run forever
	<-make(chan struct{})
}
