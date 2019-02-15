package macro

var builtins = map[string]func(string) ([]string, error){
	"fg":         fg,
	"connect":    passthrough,
	"disconnect": passthrough,
	"quit":       passthrough,

	// Internals
	"_":             passthrough,
	"_preconnect":   passthrough,
	"_connected":    passthrough,
	"_disconnected": passthrough,
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

func passthrough(args string) ([]string, error) {
	return []string{args}, nil
}
