// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

import (
	"regexp"
)

type Trigger struct {
	// The type of trigger: hilite, gag, script, macro.
	Type string

	// A regexp to match against.
	Match string

	// A list of attributes used in hilite (color, style, etc).
	Attributes []string

	// The path of a script to run.
	Script string

	// The name of a macro to run.
	Macro string

	// The compiled regexp specified in Match.
	re *regexp.Regexp
}

// compile compiles the regexp specified in the trigger's Match attribute.
func (t *Trigger) compile() error {
	re, err := regexp.Compile(t.Match)
	if err != nil {
		return err
	}
	t.re = re
	return nil
}

// Run takes the provided byte-slice from the world and, if it matches, runs
// the action specified in the trigger based on the type (hilite, gag, script
// macro). It returns the (potentially modified) input, whether or not the
// trigger matched, and any errors it encountered along the way.
func (t *Trigger) Run(input []byte) ([]byte, bool, error) {
	panic("not implemented")
	if t.re.Match(input) {
		return input, true, nil
	}
	return input, false, nil
}
