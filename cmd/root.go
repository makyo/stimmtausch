package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/makyo/stimmtausch/config"
)

var log = loggo.GetLogger("stimmtausch.cmd")

var cfgFile string

var rootCmd = &cobra.Command{
	Use:    "",
	Short:  "Run Stimmtausch.",
	Long:   "",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Stimmtausch! This is our 'coming soon' message :)")
	},
	Version: "0.0.0-pre",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home, err := config.HomeDir()
	if err != nil {
		os.Exit(2)
	}
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", filepath.Join(home, "config.yaml"), "config file")
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
