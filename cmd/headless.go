package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/makyo/stimmtausch/client"
)

func init() {
	rootCmd.AddCommand(headlessCmd)
}

var headlessCmd = &cobra.Command{
	Use:    "headless",
	Short:  "Run Stimmtausch in headless mode.",
	Long:   "",
	PreRun: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Oooh, headless! Fancy~")
		c, err := client.New()
		if err != nil {
			log.Criticalf("could not create client: %v", err)
			os.Exit(4)
		}
		log.Infof("%+v", c)
	},
	TraverseChildren: true,
}
