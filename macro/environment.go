package macro

import (
	"fmt"
	"regexp"
)

var (
	wsRE        = regexp.MustCompile("\\s+")
	macroNameRE = regexp.MustCompile("[[:alpha:]][[:word:]]*")
)

type Environment struct {
	// listeners is a list of channels to which send results from macros
	// running.
	listeners []chan MacroResult

	// macros is a map from macro name to function.
	macros map[string]func([]string) ([]string, error)
}

func (e *Environment) Dispatch(name, argString string) {
	args := wsRE.Split(argString, -1)
	var result MacroResult
	if m, ok := e.macros[name]; ok {
		results, err := m(args)
		result = MacroResult{
			Name:    name,
			Results: results,
			Err:     err,
		}
	} else {
		result = MacroResult{
			Name:    name,
			Results: args,
			Err:     fmt.Errorf("unknown macro %s", name),
		}
	}
	for _, listener := range e.listeners {
		listener <- result
	}
}

func (e *Environment) RegisterMacro(name string, m func([]string) ([]string, error)) error {
	if _, ok := e.macros[name]; ok {
		return fmt.Errorf("macro with name %s already exists", name)
	}
	if !macroNameRE.MatchString(name) {
		return fmt.Errorf("macro name must contain only letters, numbers, and underscores and start with a letter")
	}
	e.macros[name] = m
	return nil
}

func (e *Environment) AddListener(listener chan MacroResult) {
	e.listeners = append(e.listeners, listener)
}

func New() *Environment {
	return &Environment{
		macros: builtins,
	}
}
