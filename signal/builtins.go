// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package signal

import (
	"fmt"
	"strings"

	"github.com/juju/loggo"
)

var builtins = map[string]func(string) ([]string, error){
	// World switching
	"fg": fg,
	">":  func(_ string) ([]string, error) { return fg(">") },
	"<":  func(_ string) ([]string, error) { return fg("<") },

	// Connections
	"connect":    passthrough,
	"c":          passthrough,
	"disconnect": passthrough,
	"dc":         passthrough,
	"quit":       passthrough,

	// Logging
	"log": partsPassthrough,

	// Help
	"help": passthrough,

	// Reload
	"reload": passthrough,

	// Internals
	"syslog":                  syslog,
	"_":                       passthrough,
	"_util:split":             split,
	"_client:connected":       passthrough,
	"_client:disconnected":    passthrough,
	"_client:allDisconnected": passthrough,
	"_client:showModal":       titleSplit,
}

// fg handles the special case for the builtin `fg`, which sends a different
// response depending on whether switching or rotating connections.
func fg(args string) ([]string, error) {
	args = wsRE.Split(args, -1)[0]
	switch args {
	case "<":
		return []string{"rotate", "-1"}, nil
	case ">", "":
		return []string{"rotate", "1"}, nil
	default:
		return []string{"switch", args}, nil
	}
}

// syslog logs the given message at the given level. It's largely for
// debugging.
func syslog(args string) ([]string, error) {
	parts := wsRE.Split(args, 2)
	if len(parts) != 2 {
		return parts, fmt.Errorf("incorrect number of arguments; want 2 (level, log string), got %d", len(parts))
	}
	level, _ := loggo.ParseLevel(parts[0])
	log.Logf(level, "Syslog macro: %s", parts[1])
	return parts, nil
}

// titleSplit passes on the given args after splitting on the title separator,
// "::\n".
func titleSplit(args string) ([]string, error) {
	return strings.Split(args, "::\n"), nil
}

// passthrough simply passes on the given args to the listeners without taking
// any action on them.
func passthrough(args string) ([]string, error) {
	result := []string{}
	if len(args) > 0 {
		result = []string{args}
	}
	return result, nil
}

// partsPassthrough passes on the given args to the listeners, splitting them
// on whitespace to provide the string slice.
func partsPassthrough(args string) ([]string, error) {
	return wsRE.Split(args, -1), nil
}

// split passes on the given args to the listeners, splitting them on the first
// character in the string (which is not included in the result).
func split(args string) ([]string, error) {
	sep, args := args[0:1], args[1:]
	return strings.Split(args, sep), nil
}
