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
	Attributes string

	// For gags, whether or not to log the gagged string anyway.
	LogAnyway string `yaml:"log_anyway" toml:"log_anyway"`

	// The path of a script to run.
	Script string

	// Whether or not to send the output of the script or macro to the world.
	// If false, the user will be shown the output
	OutputToWorld string `yaml:"output_to_world" toml:"output_to_world"`

	// The name of a macro to run.
	Macro string

	// The compiled regexp specified in Match.
	re *regexp.Regexp
}

// compile compiles the regexp specified in the trigger's Match attribute.
func (t *Trigger) compile() error {
	switch t.Type {
	case "hilite":
	case "gag":
	case "script":
	case "macro":
		break
	default:
		return fmt.Errorf("unknown trigger type %s", t.Type)
	}
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
func (t *Trigger) Run(input string, cfg *Config) (string, error) {
	log.Tracef("running trigger %v", t)
	if matches := t.re.FindAllStringSubmatchIndex(input, -1); len(matches) != 0 {
		switch t.Type {
		case "hilite":
			return t.hiliteString(input, matches)
			break
		case "gag":
			return "", nil
			break
		case "script":
			return t.runScript(input, matches)
			break
		case "macro":
			return t.runMacro(input, matches, cfg)
			break
		}
	}
	return nil
}

// hiliteString applies ANSI escape-code highlighting to an matches within the
// provided string. It ignores submatches. If it encounters an error in the
// process, it returns as much highlighting as it got done and the error
// generated in the process.
func (t *Trigger) hiliteString(input string, matches [][]int) (string, error) {
	log.Tracef("hiliting string")
	panic("not implemented")
}

// runScript runs a script (or any executable in $PATH) with the input string
// and matches as input, each properly quoted. The matches will be sent as as
// is, in JSON format.
func (t *Trigger) runScript(input string, matches [][]int) (string, error) {
	log.Tracef("running script")
	// We could JSON-encode this, oooor...
	matchesStr := strings.Replace(fmt.Sprintf("%v", matches), " ", ",", -1)
	panic("not implemented")
}

func (t *Trigger) runMacro(input string, matches [][]int, cfg *Config) (string, error) {
	log.Tracef("running macro")
	panic("not implemented")
}
