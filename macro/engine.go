// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package macro

import (
	"bufio"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"

	"github.com/makyo/stimmtausch/signal"
)

type Engine struct {
	states   map[string]*lua.LState
	macros   map[string]*macro
	env      *signal.Dispatcher
	listener chan signal.Signal
}

func (e *Engine) AddMacro(name, filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filename)
	if err != nil {
		return err
	}
	proto, err := lua.Compile(chunk, filename)
	if err != nil {
		return err
	}
	e.macros[name] = &macro{
		name:     name,
		filename: filename,
		compiled: proto,
	}
	return nil
}

func (e *Engine) prepare() (*lua.LState, error) {
	state := lua.NewState()
	state.PreloadModule("stimmtausch", e.Loader)
	e.states[fmt.Sprintf("%+v", state)] = state
}

func (e *Engine) call(name string, arguments []string) {
	state, err := e.prepare()
	if err != nil {
		// log err
		return
	}
	defer e.CloseState(state)
	m, ok := e.macros[name]
	if !ok {
		// log err
	}
	if err := m.run(e, arguments); err != nil {
		// log err
	}
}

func (e *Engine) listen() {
	for {
		res := <-listener
		switch res.Name {
		case "call":
			e.call(res.Payload[0], res.Payload[1:])
		default:
			if res.Name[0:12] == "_macro:call:" {
				e.call(res.Name[13:], res.Payload)
			}
		}
	}
}

func (e *Engine) CloseState(state *lua.LState) {
	state.Close()
	delete(e.states, fmt.Sprintf("%+v", state))
}

func (e *Engine) Close() {
	for _, v := range e.states {
		e.CloseState(v)
	}
}

func NewEngine(env *signal.Dispatcher) (*Engine, error) {
	e := &Engine{
		states:   map[string]*lua.LState{},
		macros:   map[string]*macro{},
		env:      env,
		listener: make(chan signal.Signal),
	}
	go e.listen()
	e.env.AddListener("macro", e.listener)
	return e, nil
}
