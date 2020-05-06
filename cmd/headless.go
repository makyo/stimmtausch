// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package cmd

import (
	"os"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"

	"github.com/makyo/stimmtausch/client"
	"github.com/makyo/stimmtausch/config"
	"github.com/makyo/stimmtausch/headless"
	"github.com/makyo/stimmtausch/signal"
)

func init() {
	initFlags(headlessCmd)
	rootCmd.AddCommand(headlessCmd)
}

// headlessCmd runs Stimmtausch in headless mode, connecting to a world without
// creating the UI.
var headlessCmd = &cobra.Command{
	Use:   "headless [flags] world-or-server [world-or-server...]",
	Short: "Run Stimmtausch in headless mode (advanced).",
	Long: `Run Stimmtausch in headless mode.

This will run Stimmtausch in headless mode. That is, it will connect to any
servers or worlds you specify (see "stimmtausch help" for details on that), but not
create a user interface. Use this if you want to use your own FIFO-aware
UI such as <https://github.com/onlyhavecans/mm.vim>.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
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

		log.Criticalf("could not create client: %v", err)
		log.Tracef("creating client")
		stClient, err := client.New(cfg, env)
		if err != nil {
			log.Criticalf("could not create client: %v", err)
			os.Exit(2)
		}
		log.Tracef("created client: %+v", stClient)

		done := make(chan bool)
		h := headless.New(args, stClient)
		go h.Run(done)

		for _, arg := range args {
			go env.Dispatch("connect", arg)
		}

		<-done
	},
	TraverseChildren: true,
}
