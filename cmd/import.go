// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"

	toml "github.com/pelletier/go-toml"
	"github.com/retgits/fdio/database"
	"github.com/retgits/fdio/util"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a toml file into the database",
	Run:   runImport,
}

// init registers the command and flags
func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVar(&tomlFile, "toml", "", "The path to the TOML file (required)")
	importCmd.MarkFlagRequired("toml")
}

// runImport is the actual execution of the command
func runImport(cmd *cobra.Command, args []string) {
	// Read the file
	config, err := toml.LoadFile(tomlFile)
	if err != nil {
		log.Fatalf("Error while reading TOML content: %s\n", err.Error())
	}

	// Convert the tree to an array of maps
	arrayMap, err := util.TomlTreeToMap(config, tomlItemKey)
	if err != nil {
		log.Fatalf("Error while converting TOML to array: %s\n", err.Error())
	}

	// Get a database
	db, err := database.New(dbFile, false)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s\n", err.Error())
	}

	// Load the data into the database
	for _, item := range arrayMap {
		err = db.InsertContribution(item)
		if err != nil {
			log.Printf("Error while loading data into the database: %s\n", err.Error())
		}
	}
}
