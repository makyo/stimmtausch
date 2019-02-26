package connection

import (
	"fmt"
	"io"
	"os"

	"github.com/makyo/stimmtausch/util"
)

// output represents a named io.WriteCloser.
type output struct {
	// A name used for logging and referencing down the line.
	name string

	// Whether or not- this is the global output.
	global bool

	// The io.Writecloser itself.
	output io.WriteCloser

	// Whether or not the output supports ANSI escape codes.
	supportsANSI bool
}

// makeLogfile creates a logfile from a given name.
func (c *Connection) makeLogfile(out *output) error {
	log.Tracef("creating a log file for %s", c.name)

	log.Tracef("checking if %s exists", out.name)
	_, err := os.Stat(out.name)
	if err == nil {
		fmt.Printf("Warning: %v already exists; appending.\n", out.name)
	}

	log.Tracef("opening %s for logging", c.name)
	f, err := os.OpenFile(out.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Warningf("could not open logfile %s for %s, not logging. %v", out.name, c.name, err)
		return err
	}

	out.output = f
	log.Debugf("logfile created as %s", out.name)
	return nil
}

// closeOutputs closes open outfiles.
func (c *Connection) closeOutputs() {
	log.Tracef("closing all outputs for %s", c.name)
	for _, out := range c.outputs {
		log.Tracef("closing output file %s for %s", out.name, c.name)
		if err := out.output.Close(); err != nil {
			log.Warningf("error closing output %s for %s. %v", out.name, c.name, err)
		}
		log.Debugf("output file %s for %s closed", out.name, c.name)
		if !out.global {
			continue
		}
		if c.world.Log {
			rotateTo := c.getLogFile(fmt.Sprintf("%s.log", c.getTimestamp()))
			if err := util.StripANSIFromFile(out.name, rotateTo); err != nil {
				log.Warningf("unable to clean and rotate log file %s, you'll need to do that on your own. %v", out.name, err)
				continue
			}
		}
		log.Debugf("output file for %s rotated", c.name)
		if err := os.Remove(out.name); err != nil {
			log.Warningf("unable to remove outfile %s", out.name)
		}
	}
}

// AddOutput creates an output struct with the given io.WriteCloser. This can
// be a file, of course, but many other things as well, including the buffer
// that the UI uses.
func (c *Connection) AddOutput(name string, w io.WriteCloser, supportsANSI bool) {
	log.Tracef("creating output %s for %s", name, c.name)
	c.outputs = append(c.outputs, &output{
		name:         name,
		global:       false,
		output:       w,
		supportsANSI: supportsANSI,
	})
}
