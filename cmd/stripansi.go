package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/makyo/stimmtausch/util"
)

func init() {
	rootCmd.AddCommand(stripansiCmd)
}

var stripansiCmd = &cobra.Command{
	Use:   "strip-ansi infile outfile",
	Short: "Strip ANSI color codes from a file.",
	Long: `Strip ANSI color codes from a file.

If you have a log file which includes ANSI color codes, which might happen by
accident when a Stimmtausch session is interrupted, you can run "st strip-ansi
<input file> <output file>".`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := util.StripANSIFromFile(args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error stripping ANSI from file: %v", err)
			os.Exit(1)
		}
	},
}
