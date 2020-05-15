package macro

import (
	"github.com/yuin/gopher-lua"
)

func (e *Engine) Loader(state *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"send":   e.modSend,
		"call":   e.modCall,
		"signal": e.modSignal,
	}

	mod := state.SetFuncs(state.NewTable(), exports)
	state.Push(mod)
	return 1
}

func (e *Engine) modSend(state *lua.LState) int {
	s := signal.Signal{
		Name:    "_macro:send",
		Payload: state.ToString(1),
		Err:     nil,
	}
	e.env.DirectDispatch(s)
	return 0
}

func (e *Engine) modCall(state *lua.LState) int {
	e.call("TODO", []string{"TODO"})
	return 0
}

func (e *Engine) modSignal(state *lua.LState) int {
	s := signal.Signal{
		Name:    "TODO",
		Payload: []string{"TODO"},
		Err:     nil,
	}
	e.env.DirectDispatch(s)
	return 0
}
