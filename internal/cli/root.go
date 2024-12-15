// Package cli implements the command line interface that is used from the entrypoint.
package cli

import (
	"fmt"
	"os"

	"github.com/13excite/quic-checker/internal/cli/shell"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quic-checker",
	Short: "A CLI tool to check QUIC support of a website",
	Long:  ``,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		// log.InitCLILogger()
	},
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Usage(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(shell.NewShellCommand())
}
