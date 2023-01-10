package logon

import (
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/export"
	"github.com/vela-security/vela-public/lua"
)

var xEnv assert.Environment

/*
	local f = vela.logon.fail()
	f.ignore("xxx")
	f.pipe()
	f.start()
*/

func call(L *lua.LState) int {
	return 0
}

func WithEnv(env assert.Environment) {
	xEnv = env
	tab := lua.NewUserKV()
	tab.Set("fail", lua.NewFunction(newLogonFailL))
	tab.Set("success", lua.NewFunction(newLogonSuccessL))
	tab.Set("logout", lua.NewFunction(newLogoutL))
	ex := export.New("vela.logon.export", export.WithTable(tab), export.WithFunc(call))
	xEnv.Set("logon", ex)
}
