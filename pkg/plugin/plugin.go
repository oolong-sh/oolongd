package plugin

import (
	"fmt"
	"time"

	"github.com/oolong-sh/oolong/internal/config"
)

func InitPlugins(cfg *config.OolongConfig) {
	pm, err := NewPluginManager()
	if err != nil {
		fmt.Println("Error initializing plugin manager:", err)
		return
	}
	defer pm.LuaState.Close()

	pm.LoadPlugins(cfg.PluginPaths)
	if err != nil {
		fmt.Println("Error loading plugins:", err)
		return
	}

	err = pm.TriggerEvent("dailyNoteEvent")
	if err != nil {
		fmt.Println("Error triggering daily note event:", err)
		return
	}

	// start an event loop in a separate goroutine (act as a daemon)
	go pm.StartEventLoop()

	err = pm.TriggerEvent("customEvent", "example data")
	if err != nil {
		fmt.Println("Error triggering event:", err)
	}

	time.Sleep(5 * time.Second)
	pm.StopEventLoop()
}
