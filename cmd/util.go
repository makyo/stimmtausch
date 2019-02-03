// Stimmtausch - a MU* client - https://stimmtausch.com
//
// https://github.com/makyo/stimmtausch
// Copyright © 2019 the Stimmtausch authors
// Released under the MIT license.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"

	"github.com/makyo/st/config"
)

// initFlags constructs all of the flags that might be used by Stimmtausch.
func initFlags(cmd *cobra.Command) {
	home, err := config.HomeDir()
	if err != nil {
		os.Exit(2)
	}
	cmd.Flags().StringVarP(&cfgFile, "config", "c", filepath.Join(home, "config.yaml"), "config file")
}

// initConfig initializes the viper configuration,
func initConfig(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Criticalf("could not read in config: %v", err)
		os.Exit(3)
	}
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	loggo.ConfigureLoggers(fmt.Sprintf("<root>=%s", viper.GetString("stimmtausch.client.log_level")))
}

// GenMarkdownDocs generates markdown files for each command, which are used
// on stimmtausch.com
func GenMarkdownDocs() {
	if err := doc.GenMarkdownTree(rootCmd, "./doc/"); err != nil {
		log.Criticalf("unable to generate docs: %v", err)
		os.Exit(1)
	}
}

// GenManPages generates the man pages used for Stimmtausch.
func GenManPages() {
	header := &doc.GenManHeader{
		Title:   "STIMMTAUSCH",
		Section: "1",
	}
	if err := doc.GenManTree(rootCmd, header, "./doc/"); err != nil {
		log.Criticalf("unable to generate man pages: %v", err)
		os.Exit(1)
	}
}
