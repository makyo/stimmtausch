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

func initFlags(cmd *cobra.Command) {
	home, err := config.HomeDir()
	if err != nil {
		os.Exit(2)
	}
	cmd.Flags().StringVarP(&cfgFile, "config", "c", filepath.Join(home, "config.yaml"), "config file")
}

func initConfig(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Criticalf("could not read in config: %v", err)
		os.Exit(3)
	}
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	loggo.ConfigureLoggers(fmt.Sprintf("<root>=%s", viper.GetString("stimmtausch.client.log_level")))
}

func GenMarkdownDocs() {
	if err := doc.GenMarkdownTree(rootCmd, "./doc/"); err != nil {
		log.Criticalf("unable to generate docs: %v", err)
		os.Exit(1)
	}
}

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
