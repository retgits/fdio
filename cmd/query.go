// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"
	"os"

	"github.com/retgits/fdio/database"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Run a query against the database",
	Run:   runQuery,
}

// Flags
var (
	query string
)

// init registers the command and flags
func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&query, "query", "q", "", "The database query you want to run")
	queryCmd.MarkFlagRequired("query")
}

// runQuery is the actual execution of the command
func runQuery(cmd *cobra.Command, args []string) {
	// Get a database
	db, err := database.New(dbFile, false)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s\n", err.Error())
	}

	queryOpts := database.QueryOptions{
		Writer:     os.Stdout,
		Query:      query,
		MergeCells: true,
		RowLine:    true,
		Render:     true,
	}
	_, err = db.RunQuery(queryOpts)
	if err != nil {
		log.Printf("Error while executing query: %s\n", err.Error())
	}
}
