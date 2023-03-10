package logon

import (
	"github.com/bytedance/sonic"
	opcode "github.com/vela-security/vela-opcode"
	"github.com/vela-security/vela-public/lua"
	risk "github.com/vela-security/vela-risk"
	vtime "github.com/vela-security/vela-time"
)

func (ev *Event) String() string                         { return lua.B2S(ev.Byte()) }
func (ev *Event) Type() lua.LValueType                   { return lua.LTObject }
func (ev *Event) AssertFloat64() (float64, bool)         { return 0, false }
func (ev *Event) AssertString() (string, bool)           { return lua.B2S(ev.Byte()), true }
func (ev *Event) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (ev *Event) Peek() lua.LValue                       { return ev }

func (ev *Event) Byte() []byte {
	ev.MinionID = xEnv.ID()
	ev.Inet = xEnv.Inet()
	chunk, err := sonic.Marshal(ev)
	if err != nil {
		return nil
	}
	return chunk
}

func (ev *Event) reportL(L *lua.LState) int {
	err := xEnv.TnlSend(opcode.OpLogon, ev)
	if err != nil {
		xEnv.Debugf("logon event report fail %v data:%v", err, ev)
	}
	return 0
}

func (ev *Event) riskL(L *lua.LState) int {
	ret := risk.NewEv(risk.Class(risk.TLogin))
	ret.RemoteIP = ev.Addr
	ret.RemotePort = ev.Port
	ret.LocalIP = ev.Inet
	ret.LocalPort = ev.Port
	ret.Payload = ev.User
	ret.FromCode = L.CodeVM()
	ret.Subject = ev.Class
	L.Push(ret)
	return 1
}

func (ev *Event) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "user":
		return lua.S2L(ev.User)
	case "addr":
		return lua.S2L(ev.Addr)
	case "time":
		return vtime.VTime(ev.Time)
	case "host":
		return lua.S2L(ev.Host)
	case "pid":
		return lua.LInt(ev.Pid)
	case "class":
		return lua.S2L(ev.Class)
	case "process":
		return lua.S2L(ev.Process)
	case "type":
		return lua.S2L(ev.Typ)
	case "risk":
		return lua.NewFunction(ev.riskL)

	case "report":
		return lua.NewFunction(ev.reportL)
	}
	return lua.LNil
}
