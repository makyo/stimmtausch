// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package connection

import (
	"fmt"
	"strings"
)

func (c *Connection) parseLogSignal(args []string) error {
	if len(args) < 1 || args[0] == "" {
		args[0] = "--help"
	}
	if args[0][0:2] == "--" {
		switch args[0] {
		case "--help":
			// XXX this will trigger a help modal on all attached clients!
			// We'll probably need to come up with a way to track who dispatched
			// an event so that only _they_ can listen for it and everyone else
			// can ignore it. Make a bug for this once you have signal, Maddy.
			log.Tracef("showing help for /log")
			go c.env.Dispatch("help", "log")
		case "--off":
			c.closeLog(args[1])
		case "--list":
			c.listLogs()
		default:
			log.Warningf("unknown switch %s for log command", args[0])
		}
	} else {
		return c.openLog(args[0])
	}
	return nil
}

func (c *Connection) openLog(name string) error {
	log.Tracef("creating output %s for %s via /log", name, c.name)
	out := &output{
		name:         name,
		global:       false,
		userCreated:  true,
		supportsANSI: false,
	}
	if err := c.makeLogfile(out); err != nil {
		log.Warningf("unable to start logging. %v", err)
		return err
	}
	c.outputs = append(c.outputs, out)
	return nil
}

func (c *Connection) closeLog(name string) {
	log.Tracef("closing log %s for %s via /log", name, c.name)
	for i, out := range c.outputs {
		if out.userCreated && out.name == name {
			log.Infof("log %s closed", name)
			out.output.Close()
			c.outputs = append(c.outputs[0:i], c.outputs[i+1:]...)
			return
		}
	}
	log.Warningf("no log %s for %s", name, c.name)
}

func (c *Connection) listLogs() {
	log.Tracef("listing open logs for %s", c.name)
	logs := []string{}
	for _, out := range c.outputs {
		if out.userCreated {
			logs = append(logs, "* "+out.name)
		}
	}
	if len(logs) == 0 {
		logs = []string{"(none)"}
	}
	logList := fmt.Sprintf("Open logs for %s::\n%s", c.world.DisplayName, strings.Join(logs, "\n"))
	go c.env.Dispatch("_client:showModal", logList)
}
