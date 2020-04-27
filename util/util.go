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
	githubActivitySearchQuery = "sort=indexed&order=desc&q=filename%3Aactivity.go+flogo"
	githubTriggerSearchQuery  = "sort=indexed&order=desc&q=filename%3Atrigger.go+flogo"
)

// Crawl will search on GitHub for activity.go or trigger.go files that are related to Flogo
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
		time.Sleep(10 * time.Second)
	}
}

// prepareItems takes the raw response from the GitHub search API and turns it into an
// array of map[string]interface{} so it can be added to the database
func prepareItems(items []interface{}, db *database.Database) (map[string]interface{}, error) {
	var lastItem map[string]interface{}

	for _, item := range items {
		// For each item we need some metadata, to fillout the db
		// Map 'm' contains the path and name of the item being searched for
		// Map 'r' contains info about the repo returned
		m := item.(map[string]interface{})
		r := m["repository"].(map[string]interface{})
		author := strings.Split(r["full_name"].(string), "/")[0]

		//If the item is not a go test file , a fork, or from project-flogo, add to db and record lastItem
		if !strings.Contains(m["name"].(string), "_test") && !r["fork"].(bool) && author != "project-flogo" {
			tempMap := make(map[string]interface{})
			ref := fmt.Sprintf("https://github.com/%s/%s", r["full_name"], m["path"])
			ref = strings.Replace(ref, m["name"].(string), "", 1)
			tmplist := strings.Split(ref, "/")
			tempMap["showcase"] = ""
			c := database.Contribution{
				Ref:              ref,
				Name:             tmplist[len(tmplist)-2],
				ContributionType: strings.Replace(m["name"].(string), ".go", "", 1),
				Description:      "",
				SourceURL:        fmt.Sprintf("https://github.com/%s/tree/master/%s", r["full_name"].(string), m["path"]),
				UploadedOn:       time.Now(),
				Author:           author,
				ShowcaseEnabled:  "no",
			}
			err := db.InsertContribution(c)
			if err != nil {
				log.Printf(err.Error())
				return nil, err
			}
			lastItem = tempMap
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
