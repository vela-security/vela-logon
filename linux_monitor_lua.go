//go:build linux
// +build linux

package logon

import (
	"github.com/vela-security/vela-public/lua"
	"github.com/vela-security/vela-public/pipe"
	"time"
)

func (m *Monitor) startL(L *lua.LState) int {
	m.Start()
	m.V(lua.PTRun, time.Now())
	return 0
}

func (m *Monitor) historyL(L *lua.LState) int {
	pip := pipe.NewByLua(L)
	all := m.cat()
	if len(all) == 0 {
		return 0
	}

	for _, ev := range all {
		pip.Do(ev, m.cfg.co, func(err error) {
			xEnv.Errorf("%s history call fail %v", m.Name(), err)
		})
	}
	return 0
}

func (m *Monitor) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "start":
		return lua.NewFunction(m.startL)

	case "history":
		return lua.NewFunction(m.historyL)

	default:
		return m.cfg.Index(L, key)
	}
}
