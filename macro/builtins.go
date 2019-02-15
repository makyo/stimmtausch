package macro

import (
	"fmt"
)

var builtins = map[string]func([]string) ([]string, error){
	"_":          passthrough,
	"fg":         fg,
	"connect":    passthrough,
	"disconnect": passthrough,
}

func fg(args []string) ([]string, error) {
	if len(args) != 1 {
		return []string{}, fmt.Errorf("received wrong number of args, wanted 1, got %d (%q)", len(args), args)
	}
	switch args[0] {
	case "<":
		return []string{"rotate", "-1"}, nil
	case ">":
		return []string{"rotate", "1"}, nil
	default:
		return []string{"switch", args[0]}, nil
	}
}

func passthrough(args []string) ([]string, error) {
	return args, nil
}
