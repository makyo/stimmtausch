// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package signal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/loggo"
)

var (
	wsRE        = regexp.MustCompile("\\s+")
	macroNameRE = regexp.MustCompile("[[:alpha:]][[:word:]]*")

	log = loggo.GetLogger("stimmtausch.macro")
)

type Dispatcher struct {
	// listeners is a list of channels to which send the results of handlers
	// running.
	listeners map[string]chan Signal

	// handlers is a map from handler name to function.
	handlers map[string]func(string) ([]string, error)
}

func (e *Dispatcher) Dispatch(name, args string) {
	var result Signal
	args = strings.TrimSpace(args)
	name = strings.TrimSpace(name)
	if m, ok := e.handlers[name]; ok {
		results, err := m(args)
		result = Signal{
			Name:    name,
			Payload: results,
			Err:     err,
		}
	} else {
		// TODO here is where we'll fall back to trying macros.
		result = Signal{
			Name:    name,
			Payload: []string{args},
			Err:     fmt.Errorf("unknown macro %s", name),
		}
	}
	e.DirectDispatch(result)
}

func (e *Dispatcher) DirectDispatch(result Signal) {
	log.Tracef("dispatching %+v to %d listeners", result, len(e.listeners))
	for whence, listener := range e.listeners {
		go func(l chan Signal) { l <- result }(listener)
		log.Tracef("dispatched to %s", whence)
	}
}

func (e *Dispatcher) AddListener(whence string, listener chan Signal) {
	e.listeners[whence] = listener
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers:  builtins,
		listeners: map[string]chan Signal{},
	}
}
