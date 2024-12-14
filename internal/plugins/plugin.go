package plugins

import (
	"log"
	"time"

	"github.com/oolong-sh/oolongd/internal/config"
)

func InitPlugins() {
	pm, err := NewPluginManager()
	if err != nil {
		log.Println("Error initializing plugin manager:", err)
		return
	}
	defer pm.LuaState.Close()

	if err := pm.LoadPlugins(config.PluginPaths()); err != nil {
		log.Println("Error loading plugins:", err)
		return
	}

	if err := pm.TriggerEvent("dailyNoteEvent"); err != nil {
		log.Println("Error triggering daily note event:", err)
		return
	}

	// start an event loop in a separate goroutine (act as a daemon)
	go pm.StartEventLoop()

	if err := pm.TriggerEvent("customEvent", "example data"); err != nil {
		log.Println("Error triggering event:", err)
	}

	time.Sleep(2 * time.Second)
	pm.StopEventLoop()
}
