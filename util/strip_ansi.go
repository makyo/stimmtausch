package util

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/juju/loggo"
)

var log = loggo.GetLogger("stimmtausch.util")

// StripANSI is a regexp that matches ANSI escape codes for the purpose of
// removing them where necessary.
var StripANSI = regexp.MustCompile("\x1b\\[\\d+(;\\d+)*m")

// StripANSIFromFile reads in contents from one file line by line, strips them
// of their ANSI escape codes, then writes them to another.
func StripANSIFromFile(in, out string) error {
	log.Tracef("stripping ANSI escape codes from %s and writing to %s", in, out)
	fin, err := os.Open(in)
	if err != nil {
		return err
	}
	defer fin.Close()
	fout, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fout.Close()

	lines := 0
	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		fmt.Fprintln(fout, StripANSI.ReplaceAllString(scanner.Text(), ""))
		lines++
	}
	log.Debugf("stripped ANSI escape codes from %d lines in %s and wrote to %s", lines, in, out)
	return scanner.Err()
}
