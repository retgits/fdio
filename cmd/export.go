// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/retgits/fdio/database"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the database to a toml file",
	Run:   runExport,
}

// Flags
var (
	overwrite bool
)

// init registers the command and flags
func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&tomlFile, "toml", "", "The path to the TOML file (required)")
	exportCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite file if it exists")
	exportCmd.MarkFlagRequired("toml")
}

// runExport is the actual execution of the command
func runExport(cmd *cobra.Command, args []string) {
	// Get a database
	db, err := database.New(dbFile, false, false)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s\n", err.Error())
	}

	// Check the file
	_, err = os.Stat(tomlFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error while checking export file: %s\n", err.Error())
	} else if err == nil && !overwrite {
		// This means the file existed and the user doesn't want to overwrite
		log.Fatal("File already exists and overwrite flag is not set")
	}

	// Remove and create the file
	os.Remove(tomlFile)

	file, err := os.Create(tomlFile)
	if err != nil {
		log.Fatalf("Error while creating the new TOML file: %s\n", err.Error())
	}
	defer file.Close()

	file, err = os.OpenFile(tomlFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Error while opening TOML file: %s\n", err.Error())
	}

	// Query the database
	_, rows, err := db.DoQuery("select ref, name, type, description, url, uploadedon, author, showcase from acts")

	// Loop over the result
	for _, row := range rows {
		item := fmt.Sprintf("[[items]]\nname = %s\ntype = %s\ndescription = %s\nurl = %s\nuploadedon = %s\nauthor = %s\nshowcase = %s\n\n", row[1], row[2], row[3], row[4], row[5], row[6], row[7])
		if _, err = file.WriteString(item); err != nil {
			log.Fatalf("Error while writing %s to file: %s\n", row[0], err.Error())
		}
	}
}
