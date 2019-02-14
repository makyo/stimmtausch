// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

// Package cmd contains all of the cobra commands used to run Stimmtausch.
package cmd

import (
	"os"

	"github.com/juju/loggo"
	"github.com/spf13/cobra"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/config"
	"github.com/makyo/stimmtausch/ui"
)

var log = loggo.GetLogger("stimmtausch.cmd")

// rootCmd runs Stimmtausch with the GUI and connects to the specified world.
var rootCmd = &cobra.Command{
	Use:   "st [flags] world-or-server [world-or-server...]",
	Short: "Run Stimmtausch.",
	Long: `Run Stimmtausch.
	
Stimmtausch is a client for connecting to MU* servers. You can specify these in
your config file, along with information used to log into them. You can run the
"init" sub-command to generate a config file for you.

You may specify which worlds or servers you would like to connect to on the
command line separated by spaces. For each, Stimmtausch will first look for the
world named that in the config file, then the server named that in the config
file if no world is found. Finally, it will try to connect to that address
directly, if you specify it as "<host>:<port>".

For example, say you have a server named "furrymuck" and a world named "fm_fox".
You could connect to the world (which would be, say, the character Foxface on
FurryMUCK) with:

    st fm_fox
	
Or you could connect to FurryMUCK with:

    st furrymuck
	
Finally, if you want to connect to another server entirely, you can do so with:

    st mu.example.org:8889
	
You can combine these at will, of course. if you have the world "spr_rudder",
you could connect to bot worlds, plus a new server, with:

    st fm_fox spr_rudder mu.example.org:8889
	
For more help, see https://stimmtausch.com`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if logLevel == "" {
			initLogging("INFO")
		} else {
			initLogging(logLevel)
		}
		cfg, err := config.Load()
		if err != nil {
			log.Criticalf("unable to read config: %v", err)
			os.Exit(1)
		}
		if logLevel == "" {
			initLogging(cfg.Client.Syslog.LogLevel)
		}

		log.Tracef("creating client")
		stClient, err := client.New(cfg)
		if err != nil {
			log.Criticalf("could not create client: %v", err)
			os.Exit(2)
		}
		log.Tracef("created client: %+v", stClient)
		tui := ui.New(args, cfg, stClient)

		done := make(chan bool)
		go tui.Run(done)

		<-done
	},
	Version: "0.0.2",
}

// Execute executes the specified command via the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Criticalf("error executing the command: %v", err)
		os.Exit(2)
	}
}

// init initializes the flags required by the root command.
func init() {
	initFlags(rootCmd)
}
