// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fdio",
	Short: "Flogo Dot IO command-line",
	Long: `
A command-line interface for the Flogo Dot IO website`,
}

// Variables used in multiple flags
var (
	dbFile   string
	tomlFile string
)

const (
	tomlItemKey = "items"
	version     = "0.0.8"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbFile, "db", "", "The path to the database (required)")
	rootCmd.MarkPersistentFlagRequired("db")
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("\nYou're running FDIO version {{.Version}}\n\n")
}
