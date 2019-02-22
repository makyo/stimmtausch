package util

import (
	"os"
)

// EnsureDir ensures that a directory is in place before opening a file there.
func EnsureDir(path string) error {
	log.Tracef("ensuring %s exists", path)
	return os.MkdirAll(path, 0755)
}
