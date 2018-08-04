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
	// Get a database
	db, err := database.New(dbFile, false, false)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s\n", err.Error())
	}

	// Get the top 5 authors
	queryOpts := database.QueryOptions{
		Writer:     os.Stdout,
		Query:      "select author, count(author) as num from acts group by author order by num desc limit 5",
		MergeCells: true,
		RowLine:    true,
		Render:     true,
	}
	_, err = db.RunQuery(queryOpts)
	if err != nil {
		log.Printf("Error while executing query: %s\n", err.Error())
	}

	// Get the item types
	queryOpts.Query = "select type, count(type) as num from acts group by type"
	_, err = db.RunQuery(queryOpts)
	if err != nil {
		log.Printf("Error while executing query: %s\n", err.Error())
	}
}
