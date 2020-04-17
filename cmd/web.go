// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	initFlags(webCmd)
	rootCmd.AddCommand(webCmd)
}

// webCmd runs Stimmtausch in headless mode, connecting to a world without
// creating the UI.
var webCmd = &cobra.Command{
	Use:   "headless [flags] world-or-server [world-or-server...]",
	Short: "Run Stimmtausch in headless mode (advanced).",
	Long: `Run Stimmtausch as a webserver

This will run a websocket server which will serve MUCK connections over a 
websocket to the browser.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Criticalf("not implemented")
		os.Exit(1)
	},
	TraverseChildren: true,
}
