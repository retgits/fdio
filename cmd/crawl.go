// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/retgits/fdio/database"
	"github.com/retgits/fdio/util"
	"github.com/spf13/cobra"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawls GitHub to find new activities and triggers",
	Run:   runCrawl,
}

// Flags
var (
	actsType string
	timeout  float64
)

const (
	crawlLockFile string = ".crawl"
)

// init registers the command and flags
func init() {
	rootCmd.AddCommand(crawlCmd)
	crawlCmd.Flags().StringVar(&actsType, "type", "", "The type to look for, either trigger or activity (required)")
	crawlCmd.Flags().Float64Var(&timeout, "timeout", 0, "The number of hours between now and the last repo update")
	crawlCmd.MarkFlagRequired("type")
}

// runCrawl is the actual execution of the command
func runCrawl(cmd *cobra.Command, args []string) {
	// Creating a file with the current time
	os.Remove(crawlLockFile)
	file, err := os.OpenFile(crawlLockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("Error while opening .crawl file: %s\n", err.Error())
	}
	defer file.Close()
	currentTime := time.Now().String()
	if _, err = file.WriteString(currentTime); err != nil {
		log.Printf("Error while writing date to .crawl file: %s\n", err.Error())
	}

	// This app needs to connect to GitHub using a Personal Access Token
	githubToken, set := os.LookupEnv("GHACCESSTOKEN")
	if !set {
		log.Fatalf("GitHub Access Token is not set. Please set GHACCESSTOKEN before running this command\n")
	}

	switch strings.ToUpper(actsType) {
	case "TRIGGER":
		actsType = "Trigger"
	case "ACTIVITY":
		actsType = "Activity"
	default:
		log.Fatalf("Unknown type: %s. Please use either trigger or activity\n", actsType)
	}

	// Get a database
	db, err := database.New(dbFile, false, false)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s\n", err.Error())
	}

	// Prepare HTTP headers
	httpHeader := http.Header{"Authorization": {fmt.Sprintf("token %s", githubToken)}}

	err = util.Crawl(httpHeader, db, timeout, actsType)
	if err != nil {
		log.Fatalf("Error while crawling for %s: %s\n", actsType, err.Error())
	}
	log.Printf("Completed crawling for %s!\n", actsType)
}
