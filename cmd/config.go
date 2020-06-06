package cmd

import (
	"fmt"
	"os"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/spf13/cobra"

	"github.com/makyo/stimmtausch/config"
)

var defaultConfigOnly bool

func init() {
	configCmd.Flags().BoolVar(&defaultConfigOnly, "default", false, "only dump the default configuration")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Retrieve information about configuration",
	Long: `Retrieve information about configuration
	
This command loads all configuration as Stimmtausch will see it and then dumps
it to output so that you can see the eventual result.`,
	Run: func(cmd *cobra.Command, args []string) {
		loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
		if logLevel == "" {
			initLogging("INFO")
		} else {
			initLogging(logLevel)
		}
		if defaultConfigOnly {
			fmt.Println(config.DefaultConfig)
			return
		}
		cfg, err := config.New()
		if err != nil {
			log.Criticalf("unable to read config: %v", err)
			os.Exit(1)
		}
		fmt.Println(cfg.Dump())
	},
}
