// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/makyo/st/config"
	"github.com/makyo/st/headless"
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
servers or worlds you specify (see "st help" for details on that), but not
create a user interface. Use this if you want to use your own FIFO-aware
UI such as <https://github.com/onlyhavecans/mm.vim>.`,
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
		headless.New(args, cfg)
	},
	TraverseChildren: true,
}
