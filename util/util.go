// Package util implements utility methods
package util

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/retgits/fdio/database"
	"github.com/tomnomnom/linkheader"
)

const (
	githubRootEndpoint        = "https://api.github.com"
	githubSearchEndpoint      = "/search/code"
	githubActivitySearchQuery = "sort=indexed&order=desc&q=filename%3Aactivity.json+flogo"
	githubTriggerSearchQuery  = "sort=indexed&order=desc&q=filename%3Atrigger.json+flogo"
)

// Crawl will search on GitHub for activity.json or trigger.json files that are related to Flogo
func Crawl(httpHeader http.Header, db *database.Database, timeout float64, contribType string) error {
	githubSearchQuery := ""
	if contribType == "Trigger" {
		githubSearchQuery = githubTriggerSearchQuery
	} else {
		githubSearchQuery = githubActivitySearchQuery
	}

	var maxPages int
	i := 0

	for {
		// Get a GitHub page
		URL := fmt.Sprintf("%s%s?%s&page=%v", githubRootEndpoint, githubSearchEndpoint, githubSearchQuery, i)
		log.Printf("Send request to %s", URL)
		response, err := HTTPRequest(URL, httpHeader)
		if err != nil {
			return err
		}

		// Check how many pages exist
		// Only do this the first time
		if i == 0 {
			links := linkheader.Parse(response.Headers.Get("Link"))
			for _, link := range links {
				if link.Rel == "last" {
					maxPages, _ = strconv.Atoi(link.URL[strings.Index(link.URL, "page=")+5:])
				}
			}
			log.Printf("There are a total of %d pages", maxPages)
		}

		// Collect the items in this set
		lastItem, err := prepareItems(response.Body["items"].([]interface{}), db)
		if err != nil {
			return fmt.Errorf("error while converting response to array: %s", err.Error())
		}

		// Check the last update time
		duration, err := checkLastRepoUpdate(lastItem, httpHeader)
		// If update is larger than timeout it means the last update to the last checked
		// repository was longer than the timeout we set. In that case we don't need to
		// scan any further
		if duration > timeout {
			log.Printf("Maximum timeout reached. Last repo update was %v hours\n", duration)
			return nil
		}

		// Stop if the maxPages is reached
		if i++; i == maxPages {
			return nil
		}

		// Wait for 5 seconds so the GitHub search API limit won't be breached
		time.Sleep(5 * time.Second)
	}
}

// prepareItems takes the raw response from the GitHub search API and turns it into an
// array of map[string]interface{} so it can be added to the database
func prepareItems(items []interface{}, db *database.Database) (map[string]interface{}, error) {
	var lastItem map[string]interface{}

	for _, item := range items {
		// For each item we need some metadata
		// the full_name of the repository contains both the repository name as
		// well as the name of the owner
		project := item.(map[string]interface{})
		repository := project["repository"].(map[string]interface{})
		contentURL := strings.Replace(project["html_url"].(string), "github.com", "raw.githubusercontent.com", 1)
		contentURL = strings.Replace(contentURL, "/blob", "", 1)

		// Get the content of the actual file
		response, err := HTTPRequest(contentURL, nil)

		if err == nil && response.Body["type"] != nil {
			// Get the project path without activity.json at the end
			projectPath := project["path"].(string)
			if strings.Contains(projectPath, "activity.json") {
				projectPath = projectPath[:len(projectPath)-13]
			} else {
				projectPath = projectPath[:len(projectPath)-12]
			}

			// Get the project type
			projectType := response.Body["type"].(string)
			projectType = projectType[6:]
			projectType = strings.Replace(projectType, ":", "", -1)

			// Set author to unknown if it doesn't exist
			if response.Body["author"] == nil {
				response.Body["author"] = "Unknown"
			}

			// Set an empty string if the description doesn't exist
			if response.Body["description"] == nil {
				response.Body["description"] = ""
			}

			if response.Body["name"] != nil && response.Body["ref"] != nil {
				tempMap := make(map[string]interface{})
				tempMap["ref"] = response.Body["ref"].(string)
				tempMap["name"] = response.Body["name"].(string)
				tempMap["type"] = projectType
				tempMap["description"] = response.Body["description"].(string)
				tempMap["url"] = fmt.Sprintf("https://github.com/%s/tree/master/%s", repository["full_name"].(string), projectPath)
				tempMap["uploadedon"] = ""
				tempMap["author"] = response.Body["author"].(string)
				tempMap["showcase"] = ""
				err := db.InsertContribution(tempMap)
				if err != nil {
					log.Printf(err.Error())
					return nil, err
				}
				lastItem = tempMap
			} else {
				log.Printf("%s has no name or ref field so cannot be added to FDIO", fmt.Sprintf("https://github.com/%s/tree/master/%s", repository["full_name"].(string), projectPath))
			}
		}
	}
	return lastItem, nil
}

func checkLastRepoUpdate(item map[string]interface{}, httpHeader http.Header) (float64, error) {
	url := item["url"].(string)
	urlItems := strings.Split(url, "/")

	response, err := HTTPRequest(fmt.Sprintf("%s/repos/%s/%s", githubRootEndpoint, urlItems[3], urlItems[4]), httpHeader)
	if err != nil {
		log.Print(err.Error())
		return 0, err
	}

	updatedAt := response.Body["updated_at"].(string)
	layout := "2006-01-02T15:04:05Z"
	t, _ := time.Parse(layout, updatedAt)
	duration := time.Since(t)

	return duration.Hours(), nil
}
