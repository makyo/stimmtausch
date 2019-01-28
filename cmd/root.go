package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "",
	Short: "",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to stimmtausch")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
