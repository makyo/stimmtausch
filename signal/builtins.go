package signal

import (
	"fmt"

	"github.com/juju/loggo"
)

var builtins = map[string]func(string) ([]string, error){
	"fg":         fg,
	"connect":    passthrough,
	"disconnect": passthrough,
	"quit":       passthrough,
	"syslog":     syslog,

	// Internals
	"_":                       passthrough,
	"_client:connected":       passthrough,
	"_client:disconnected":    passthrough,
	"_client:allDisconnected": passthrough,
}

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

func syslog(args string) ([]string, error) {
	parts := wsRE.Split(args, 2)
	if len(parts) != 2 {
		return parts, fmt.Errorf("incorrect number of arguments; want 2 (level, log string), got %d", len(parts))
	}
	level, _ := loggo.ParseLevel(parts[0])
	log.Logf(level, "Syslog macro: %s", parts[1])
	return parts, nil
}

func passthrough(args string) ([]string, error) {
	return []string{args}, nil
}
