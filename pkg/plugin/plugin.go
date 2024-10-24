package plugin

import (
	"errors"
	"fmt"
	"time"

	"github.com/Shopify/go-lua"
)

// TODO:
type Plugin struct {
}

func AddPlugin() error {
	// TODO:
	return nil
}

// CHANGE: remove this later, replace with a better implementation
func LuaPlugin() error {
	l := lua.NewState()
	lua.OpenLibraries(l)

	if err := lua.DoFile(l, "./scripts/daily-note.lua"); err != nil {
		return err
	}
	l.Global("CreateDailyNote")
	if !l.IsFunction(-1) {
		return errors.New("CreateDailyNote is not a function")
	}
	date := time.Now().Format("01-02-06")
	l.PushString(date)
	if err := l.ProtectedCall(1, 1, 0); err != nil {
		return err
	}
	res, ok := l.ToString(-1)
	if !ok {
		return errors.New("Failed to convert result to string.")
	}

	fmt.Println(res)

	return nil
}
