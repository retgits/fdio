// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"
	"os"

	"github.com/retgits/fdio/database"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get statistics from the database",
	Run:   runGetStats,
}

// init registers the command and flags
func init() {
	rootCmd.AddCommand(statsCmd)
}

// runGetStats is the actual execution of the command
func runGetStats(cmd *cobra.Command, args []string) {
	db := database.MustOpenSession(databaseFile)

	for _, q := range statisticsQueries {
		queryOpts := database.QueryOptions{
			Writer:     os.Stdout,
			Query:      q,
			MergeCells: true,
			RowLine:    true,
			Render:     true,
		}
		_, err := db.Query(queryOpts)
		if err != nil {
			log.Fatalf("Error while executing query: %s\n", err.Error())
		}
	}
}
