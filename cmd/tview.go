// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

// Package cmd contains all of the cobra commands used to run Stimmtausch.
package cmd

import (
	"os"

	"github.com/pkg/profile"
	"github.com/spf13/cobra"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/config"
	"github.com/makyo/stimmtausch/signal"
	ui "github.com/makyo/stimmtausch/ui/tview"
)

// rootCmd runs Stimmtausch with the GUI and connects to the specified world.
var tviewCmd = &cobra.Command{
	Use:   "tview [flags] [world-or-server...]",
	Short: "Run Stimmtausch uising the tview UI.",
	Long: `Run Stimmtausch using the tview UI.
	
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
	Args: cobra.ArbitraryArgs,
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

		if cfg.Client.Profile.CPU {
			defer profile.Start().Stop()
		} else if cfg.Client.Profile.Mem {
			defer profile.Start(profile.MemProfile).Stop()
		}

		if logLevel == "" {
			initLogging(cfg.Client.Syslog.LogLevel)
		}

		env := signal.NewDispatcher()

		log.Tracef("creating client")
		stClient, err := client.New(cfg, env)
		if err != nil {
			log.Criticalf("could not create client: %v", err)
			os.Exit(2)
		}
		log.Tracef("created client: %+v", stClient)

		done := make(chan bool)
		tview, err := ui.New(stClient)
		if err != nil {
			log.Criticalf("could not create UI: %v", err)
			os.Exit(2)
		}
		go tview.Run(done)

		for _, arg := range args {
			go env.Dispatch("connect", arg)
		}

		<-done
	},
}

// init initializes the flags required by the root command.
func init() {
	initFlags(tviewCmd)
	rootCmd.AddCommand(tviewCmd)
}
