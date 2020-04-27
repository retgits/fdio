// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"

	"github.com/retgits/fdio/database"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database in a new location",
	Run:   runInit,
}

// init registers the command and flags
func init() {
	rootCmd.AddCommand(initCmd)
}

// runInit is the actual execution of the command
func runInit(cmd *cobra.Command, args []string) {
	err := database.MustOpenSession(dbFile).Initialize()
	if err != nil {
		log.Println(err.Error())
	}
}
