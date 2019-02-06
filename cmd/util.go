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

	"github.com/makyo/st/config"
)

var (
	cfgFile  string
	logLevel string
)

// initFlags constructs all of the flags that might be used by Stimmtausch.
func initFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cfgFile, "config", "c", filepath.Join(config.ConfigDir, "config.yaml"), "config file")
	cmd.Flags().StringVarP(&logLevel, "log-level", "", "", "level of detail to show in logs (can be TRACE, DEBUG, INFO, WARNING, ERROR, CRITICAL)")
}

// initConfig initializes the configuration used within the session.
func initLogging(logLevel string) {
	// TODO log to file
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	loggo.ConfigureLoggers(fmt.Sprintf("<root>=%s", logLevel))
}

// GenMarkdownDocs generates markdown files for each command, which are used
// on stimmtausch.com
func GenMarkdownDocs() {
	if err := doc.GenMarkdownTree(rootCmd, "./doc/"); err != nil {
		log.Criticalf("unable to generate docs: %v", err)
		os.Exit(2)
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
		os.Exit(2)
	}
}
