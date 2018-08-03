// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version",
	Run:   runGetVersion,
}

// init registers the command and flags
func init() {
	rootCmd.AddCommand(versionCmd)
}

// runGetVersion is the actual execution of the command
func runGetVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("FDIO version: %s\n", version)
}
