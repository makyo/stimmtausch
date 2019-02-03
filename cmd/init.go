// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/makyo/st/config"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:              "init",
	Short:            "Initialize a config file if none exists.",
	Long:             "If no default configuration file exists, initialize one with some sensible defaults plus some hints of what you'll need to fill in.",
	Run:              config.Initialize,
	TraverseChildren: true,
}
