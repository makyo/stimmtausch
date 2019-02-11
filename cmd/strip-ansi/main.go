package main

import (
	"fmt"
	"os"

	"github.com/makyo/stimmtausch/util"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Strip ANSI escape codes from a file.\n\n    USAGE: %s infile outfile\n\n", os.Args[0])
		os.Exit(1)
	}
	if err := util.StripANSIFromFile(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "Error stripping ANSI from file: %v", err)
		os.Exit(1)
	}
}
