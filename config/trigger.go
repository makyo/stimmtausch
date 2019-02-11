// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package config

import (
	"fmt"
	"regexp"
	"strings"

	ansi "github.com/makyo/ansigo"
)

type Trigger struct {
	// The type of trigger: hilite, gag, script, macro.
	Type string

	// A regexp to match against.
	Match string

	// A list of regexps to match against.
	Matches []string

	// A list of attributes used in hilite (color, style, etc).
	Attributes string

	// For gags, whether or not to log the gagged string anyway.
	LogAnyway bool `yaml:"log_anyway" toml:"log_anyway"`

	// The path of a script to run.
	Script string

	// Whether or not to send the output of the script or macro to the world.
	// If false, the user will be shown the output
	OutputToWorld string `yaml:"output_to_world" toml:"output_to_world"`

	// The name of a macro to run.
	Macro string

	// The compiled regexp specified in Match.
	reList []*regexp.Regexp
}

// compile compiles the regexp specified in the trigger's Match attribute.
func compileTrigger(t Trigger) (*Trigger, error) {
	switch t.Type {
	case "hilite":
	case "gag":
	case "script":
	case "macro":
		break
	default:
		return nil, fmt.Errorf("unknown trigger type %s", t.Type)
	}
	if t.Match == "" && len(t.Matches) == 0 {
		return nil, fmt.Errorf("no matches for trigger")
	}
	if t.Match != "" {
		re, err := regexp.Compile(t.Match)
		if err != nil {
			return nil, err
		}
		t.reList = append(t.reList, re)
	}
	for _, match := range t.Matches {
		re, err := regexp.Compile(match)
		if err != nil {
			return nil, err
		}
		t.reList = append(t.reList, re)
	}
	return &t, nil
}

// Run takes the provided byte-slice from the world and, if it matches, runs
// the action specified in the trigger based on the type (hilite, gag, script
// macro). It returns the (potentially modified) input, whether or not the
// trigger matched, and any errors it encountered along the way.
func (t *Trigger) Run(input string, cfg *Config) (bool, string, []error) {
	log.Tracef("running trigger %+v", t)
	applies := false
	var errs []error
	for _, re := range t.reList {
		if matches := re.FindAllStringSubmatchIndex(input, -1); len(matches) != 0 {
			applies = true
			var err error
			switch t.Type {
			case "hilite":
				input, err = t.hiliteString(input, matches)
				break
			case "gag":
				return true, input, []error{}
				break
			case "script":
				input, err = t.runScript(input, matches)
				break
			case "macro":
				input, err = t.runMacro(input, matches, cfg)
				break
			}
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return applies, input, errs
}

// hiliteString applies ANSI escape-code highlighting to an matches within the
// provided string. It ignores submatches. If it encounters an error in the
// process, it returns as much highlighting as it got done and the error
// generated in the process.
func (t *Trigger) hiliteString(input string, matches [][]int) (string, error) {
	log.Tracef("hiliting string")
	var parts []string
	offset := 0
	for _, match := range matches {
		before, target, after := input[:match[0]-offset], input[match[0]-offset:match[1]-offset], input[match[1]-offset:]
		offset = match[1]
		// We need to use ApplyWithReset here because termbox doesn't support
		// color-off codes
		target, err := ansi.Apply(t.Attributes, target)
		if err != nil {
			log.Warningf("error applying hilites: %v (continuing anyway)", err)
		}
		parts = append(parts, before, target)
		input = after
	}
	parts = append(parts, input)
	return strings.Join(parts, ""), nil
}

// runScript runs a script (or any executable in $PATH) with the input string
// and matches as input, each properly quoted. The matches will be sent as as
// is, in JSON format.
func (t *Trigger) runScript(input string, matches [][]int) (string, error) {
	log.Tracef("running script")
	// We could JSON-encode this, oooor...
	matchesStr := strings.Replace(fmt.Sprintf("%v", matches), " ", ",", -1)
	matchesStr = matchesStr
	return input, fmt.Errorf("not implemented")
}

func (t *Trigger) runMacro(input string, matches [][]int, cfg *Config) (string, error) {
	log.Tracef("running macro")
	return input, fmt.Errorf("not implemented")
}
