package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	//
	//"github.com/makyo/stimmtausch/client"
)

func init() {
	rootCmd.AddCommand(headlessCmd)
}

var headlessCmd = &cobra.Command{
	Use:   "headless",
	Short: "Run Stimmtausch in headless mode.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Oooh, headless! Fancy~")
	},
}
