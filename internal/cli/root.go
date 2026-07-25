package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
	format  string
)

var rootCmd = &cobra.Command{
	Use:   "netscope",
	Short: "Reconnaissance and security analysis platform",
	Long:  "NetScope discovers assets, analyzes services, identifies findings, and maps relationships in authorized environments.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().StringVarP(&format, "output", "o", "text", "output format: text, json, or csv")
}

func Execute() error {
	return rootCmd.Execute()
}

func printErr(msg string) {
	if !quiet {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", msg)
	}
}
