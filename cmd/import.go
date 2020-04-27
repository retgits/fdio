// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"log"
	"time"

	toml "github.com/pelletier/go-toml"
	"github.com/retgits/fdio/database"
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
	arrayMap, err := TomlTreeToContributions(config, tomlItemKey)
	if err != nil {
		log.Fatalf("Error while converting TOML to array: %s\n", err.Error())
	}

	// Get a database
	db, err := database.OpenSession(dbFile)
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

// TomlTreeToContributions converts a toml tree to an array of contributions. It does so
// by introspecting the tree and looking for the items that match a specific key.
func TomlTreeToContributions(tree *toml.Tree, key string) ([]database.Contribution, error) {
	// Get the correct key
	queryResult := tree.Get(key)
	if queryResult == nil {
		return nil, fmt.Errorf("No items found in the tree")
	}

	// Prepare the result
	resultArray := queryResult.([]*toml.Tree)
	datamap := make([]database.Contribution, len(resultArray))
	for idx, val := range resultArray {
		o := val.ToMap()
		datamap[idx] = database.Contribution{
			Author:           o["author"].(string),
			ContributionType: o["type"].(string),
			Description:      o["description"].(string),
			Homepage:         o["homepage"].(string),
			Name:             o["name"].(string),
			Ref:              o["ref"].(string),
			ShowcaseEnabled:  o["showcase"].(string),
			SourceURL:        o["url"].(string),
			Title:            o["title"].(string),
			UploadedOn:       time.Now(),
			Version:          o["version"].(string),
		}
	}
	return datamap, nil
}
