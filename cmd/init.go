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

// Flags
var (
	dbCreate bool
	dbReset  bool
)

// init registers the command and flags
func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&dbCreate, "create", false, "Create a new database if file doesn't exist")
	initCmd.Flags().BoolVar(&dbReset, "reset", false, "Reset the database")
}

// runInit is the actual execution of the command
func runInit(cmd *cobra.Command, args []string) {
	_, err := database.New(dbFile, dbCreate, dbReset)
	if err != nil {
		log.Printf(err.Error())
	}
}
