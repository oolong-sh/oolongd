package plugins

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type PluginManager struct {
	LuaState *lua.LState
	quit     chan bool
}

func NewPluginManager() (*PluginManager, error) {
	l := lua.NewState()
	if err := l.DoFile("./scripts/plugin_manager.lua"); err != nil {
		return nil, fmt.Errorf("error loading plugin manager: %v", err)
	}

	return &PluginManager{
		LuaState: l,
		quit:     make(chan bool),
	}, nil
}

func (pm *PluginManager) LoadPlugins(plugins []string) error {
	args := []lua.LValue{}
	for _, p := range plugins {
		args = append(args, toLuaValue(pm.LuaState, p))
	}

	err := pm.LuaState.CallByParam(lua.P{
		Fn:      pm.LuaState.GetGlobal("LoadPlugins"),
		NRet:    0,
		Protect: true,
	}, args...)
	if err != nil {
		return fmt.Errorf("error loading all plugins: %v", err)
	}
	return nil
}

func (pm *PluginManager) TriggerEvent(eventName string, args ...interface{}) error {
	l := pm.LuaState

	// prepare lua arguments (eventName, args...)
	luaArgs := []lua.LValue{lua.LString(eventName)}
	for _, arg := range args {
		luaArgs = append(luaArgs, toLuaValue(l, arg))
	}

	// call lua 'TriggerEvent' function
	err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("TriggerEvent"),
		NRet:    0,
		Protect: true,
	}, luaArgs...)
	if err != nil {
		return fmt.Errorf("error triggering event %s: %v", eventName, err)
	}

	return nil
}

func (pm *PluginManager) StartEventLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.processTimers()
		case <-pm.quit:
			return
		}
	}
}

func (pm *PluginManager) StopEventLoop() {
	pm.quit <- true
}

func (pm *PluginManager) processTimers() error {
	return pm.LuaState.CallByParam(lua.P{
		Fn:      pm.LuaState.GetGlobal("ProcessTimers"),
		NRet:    0,
		Protect: true,
	})
}

func toLuaValue(l *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	default:
		return lua.LNil
	}
}
