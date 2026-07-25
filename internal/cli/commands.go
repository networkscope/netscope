package cli

import (
	"fmt"
	"os"

	"github.com/networkscope/netscope/internal/api"
	"github.com/networkscope/netscope/internal/changes"
	"github.com/networkscope/netscope/internal/core"
	"github.com/spf13/cobra"
)

var savePath string

func init() {
	scanCmd.Flags().StringVarP(&savePath, "save", "s", "", "save assessment to SQLite database")
}

var scanCmd = &cobra.Command{
	Use:   "scan <target>",
	Short: "Scan a target for assets and services",
	Long:  "Perform reconnaissance on a target host, domain, or IP address.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		engine := core.NewEngine()
		_, err := engine.Scan(args[0])
		if err != nil {
			printErr(err.Error())
			os.Exit(1)
		}
		if savePath != "" {
			if err := engine.Save(savePath); err != nil {
				printErr("save failed: " + err.Error())
				os.Exit(1)
			}
			if !quiet {
				fmt.Fprintf(os.Stderr, "Saved assessment to %s\n", savePath)
			}
		}
		printResults(engine, args[0])
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve <address>",
	Short: "Start the API server",
	Long:  "Launch the NetScope HTTP API server for dashboard integration.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := api.NewServer(args[0])
		fmt.Fprintf(os.Stderr, "NetScope API listening on %s\n", args[0])
		if err := server.Start(); err != nil {
			printErr("server failed: " + err.Error())
			os.Exit(1)
		}
	},
}

var changesCmd = &cobra.Command{
	Use:   "changes <path>",
	Short: "Show changes since last saved assessment",
	Long:  "Compare the current in-memory state against the latest snapshot in the database.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		engine := core.NewEngine()
		if err := engine.Load(args[0]); err != nil {
			printErr("load failed: " + err.Error())
			os.Exit(1)
		}
		prev, err := engine.LoadPreviousSnapshot(args[0])
		if err != nil {
			printErr("load previous snapshot failed: " + err.Error())
			os.Exit(1)
		}
		current, err := engine.Snapshot("")
		if err != nil {
			printErr("snapshot failed: " + err.Error())
			os.Exit(1)
		}
		result := changes.Diff(prev, current)
		printChanges(result)
	},
}

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "List discovered assets",
	Long:  "Display all discovered assets from the current assessment.",
	Run: func(cmd *cobra.Command, args []string) {
		printErr("assets requires an active scan context; use 'netscope scan <target>'")
	},
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List discovered services",
	Long:  "Display all discovered services from the current assessment.",
	Run: func(cmd *cobra.Command, args []string) {
		printErr("services requires an active scan context; use 'netscope scan <target>'")
	},
}

var findingsCmd = &cobra.Command{
	Use:   "findings",
	Short: "List security findings",
	Long:  "Display all security findings from the current assessment.",
	Run: func(cmd *cobra.Command, args []string) {
		printErr("findings requires an active scan context; use 'netscope scan <target>'")
	},
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Display relationship graph",
	Long:  "Show the asset and service relationship model for the current assessment.",
	Run: func(cmd *cobra.Command, args []string) {
		printErr("graph requires an active scan context; use 'netscope scan <target>'")
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate assessment report",
	Long:  "Generate a structured report for the current assessment.",
	Run: func(cmd *cobra.Command, args []string) {
		printErr("report is not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(changesCmd)
	rootCmd.AddCommand(assetsCmd)
	rootCmd.AddCommand(servicesCmd)
	rootCmd.AddCommand(findingsCmd)
	rootCmd.AddCommand(graphCmd)
	rootCmd.AddCommand(reportCmd)
}
