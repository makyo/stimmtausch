package signal

import (
	"fmt"

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

	// Internals
	"syslog":                  syslog,
	"_":                       passthrough,
	"_client:connected":       passthrough,
	"_client:disconnected":    passthrough,
	"_client:allDisconnected": passthrough,
}

// fg handles the special case for the builtin `fg`, which sends a different
// response depending on whether switching or rotating connections.
func fg(args string) ([]string, error) {
	args = wsRE.Split(args, -1)[0]
	switch args {
	case "<":
		return []string{"rotate", "-1"}, nil
	case ">":
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
