// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/retgits/fdio/database"
	"github.com/retgits/fdio/github"
	"github.com/spf13/cobra"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawls GitHub to find new activities and triggers",
	Run:   runCrawl,
}

// init registers the command and flags
func init() {
	rootCmd.AddCommand(crawlCmd)
	crawlCmd.Flags().StringVar(&activityType, "type", "", "The type to look for: trigger, activity, or contribution (required)")
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
	githubToken, set := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !set {
		log.Fatalf("GitHub Access Token is not set. Please set GITHUB_ACCESS_TOKEN before running this command\n")
	}

	var contributionType github.ContributionIdentifier

	switch strings.ToUpper(activityType) {
	case github.TriggerType.String():
		contributionType = github.TriggerType
	case github.ActivityType.String():
		contributionType = github.ActivityType
	case github.ContributionType.String():
		contributionType = github.ContributionType
	default:
		log.Fatalf("Unknown type: %s. Please use either trigger or activity\n", activityType)
	}

	// Get a database
	db := database.MustOpenSession(databaseFile)

	err = github.Crawl(githubToken, db, timeout, contributionType)
	if err != nil {
		log.Fatalf("Error while crawling for %s: %s\n", activityType, err.Error())
	}
	log.Printf("Completed crawling for %s!\n", activityType)
}
