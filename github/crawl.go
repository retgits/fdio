package github

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/retgits/fdio/database"
)

const (
	apiEndpoint       = "https://api.github.com"
	searchPath        = "search/code"
	activityQuery     = "sort=indexed&order=desc&q=filename%3Aactivity.json+flogo"
	triggerQuery      = "sort=indexed&order=desc&q=filename%3Atrigger.json+flogo"
	contributionQuery = "sort=indexed&order=desc&q=filename%3Adescriptor.json+flogo"
)

type ContributionIdentifier int

const (
	ActivityType ContributionIdentifier = iota
	ContributionType
	TriggerType
)

func (c ContributionIdentifier) String() string {
	return [...]string{
		"ACTIVITY",
		"TRIGGER",
		"CONTRIBUTION",
	}[c]
}

// Crawl will search on GitHub for files that are related to Flogo
func Crawl(token string, db *database.Database, timeout float64, ci ContributionIdentifier) error {
	var searchQuery string
	var legacy bool
	var pathString string

	switch ci {
	case TriggerType:
		searchQuery = triggerQuery
		legacy = true
		pathString = "trigger.json"
	case ActivityType:
		searchQuery = activityQuery
		legacy = true
		pathString = "activity.json"
	case ContributionType:
		searchQuery = contributionQuery
		legacy = false
		pathString = "descriptor.json"
	}

	var maxPages int
	i := 0

	for {
		// Prepare URL
		URL := fmt.Sprintf("%s/%s?%s&page=%v", apiEndpoint, searchPath, searchQuery, i)

		res, err := getSearchResults(URL, token)
		if err != nil {
			return err
		}

		// Check how many pages exist
		// Only do this the first time
		if i == 0 {
			maxPages = getMaxPages(res.HTTPHeaders)
		}

		// Add the items to the database
		for _, repo := range res.Items {
			activityURL := strings.Replace(repo.HTMLURL, "github.com", "raw.githubusercontent.com", 1)
			activityURL = strings.ReplaceAll(activityURL, "blob/", "")

			activity, err := getActivityContent(activityURL)
			if err != nil {
				log.Printf("unable to get data for %s: %s", repo.HTMLURL, err.Error())
				continue
			}

			path := strings.Replace(repo.Path, pathString, "", 1)

			contribution := database.Contribution{
				Author:           repo.Repository.Owner.Login,
				ContributionType: ci.String(),
				Description:      activity.Description,
				Homepage:         activity.Homepage,
				Legacy:           legacy,
				Name:             activity.Name,
				Ref:              activity.Ref,
				ShowcaseEnabled:  false,
				SourceURL:        fmt.Sprintf("https://github.com/%s/tree/master/%s", repo.Repository.FullName, path),
				Title:            activity.Title,
				UploadedOn:       time.Now().Format("2006-01-02"),
				Version:          activity.Version,
			}

			err = db.InsertContribution(contribution)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed: contributions.sourceurl") {
					err = db.UpdateContribution(contribution)
					if err != nil {
						log.Printf("unable to update data for %s (%s): %s", activity.Title, repo.Repository.FullName, err.Error())
						continue
					}
				} else {
					log.Printf("unable to add %s (%s) to database: %s", activity.Title, repo.Repository.FullName, err.Error())
					continue
				}
			}

			log.Printf("added %s (%s) to database", activity.Title, repo.Repository.FullName)
		}

		// Check the last update time
		lastActivity := res.Items[len(res.Items)-1]
		duration, err := repoLastUpdated(lastActivity.Repository.FullName)
		if err != nil {
			log.Printf("unable to determine last update of %s to database: %s", lastActivity.Repository.FullName, err.Error())
		}

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

func repoLastUpdated(repo string) (float64, error) {
	url := fmt.Sprintf("%s/repos/%s", apiEndpoint, repo)

	res, err := getRepoDetails(url)
	if err != nil {
		return 0, err
	}

	layout := "2006-01-02T15:04:05Z"
	t, _ := time.Parse(layout, res.UpdatedAt)
	duration := time.Since(t)

	return duration.Hours(), nil
}
