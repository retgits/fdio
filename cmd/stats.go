// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
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

	// Prepare a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Statistic", "Description"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	// Get the top 5 authors
	results, err := db.DoStatsQuery("select author, count(author) as num from acts group by author order by num desc limit 5")
	if err != nil {
		log.Printf("Error while executing query: %s\n", err.Error())
	}

	// Add the top 5 authors to the table
	for _, item := range results {
		table.Append([]string{"Top 5 authors", item})
	}

	// Get the type counts
	results, err = db.DoStatsQuery("select type, count(type) as num from acts group by type")
	if err != nil {
		log.Printf("Error while executing query: %s\n", err.Error())
	}

	// Add the top 5 authors to the table
	for _, item := range results {
		table.Append([]string{"Type", item})
	}

	// Print the table
	table.Render()
}
