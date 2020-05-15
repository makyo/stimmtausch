// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package macro

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

type macro struct {
	name     string
	filename string
	compiled *lua.FunctionProto
}

func (m *Macro) run(state *lua.LState, arguments []string) error {
	argsTable := state.NewTable()
	for i, arg := range arguments {
		argsTable.Append(lua.LString(arg))
	}
	state.SetGlobal("arguments", argsTable)
	lfunc := state.NewFunctionFromProto(m.compiled)
	state.Push(lfunc)
	return state.PCall(0, lua.MultRet, nil)
}
