// Package cli implements the command line interface that is used from the entrypoint.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quic-checker",
	Short: "A CLI tool to check QUIC support of a website",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// log.InitCLILogger()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
