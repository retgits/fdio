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

// Flags
var (
	databaseFile string
	activityType string
	timeout      float64
)

// Queries
var (
	statisticsQueries = []string{
		"select author, count(author) as num from contributions group by author order by num desc limit 5",
		"select type, count(type) as num from contributions group by type",
	}
)

const (
	// Name of the lock file to prevent two instances of FDIO accessing resources at the same time
	crawlLockFile = ".crawl"

	// Version number of FDIO
	Version = "0.1.2"
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
	rootCmd.PersistentFlags().StringVar(&databaseFile, "db", "", "The path to the database (required)")
	rootCmd.MarkPersistentFlagRequired("db")
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("\nYou're running FDIO version {{.Version}}\n\n")
}
