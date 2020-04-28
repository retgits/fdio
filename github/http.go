package github

// The imports
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tomnomnom/linkheader"
)

// GithubData contains repositories with data and HTTP headers
type GithubData struct {
	Items       []Item
	HTTPHeaders http.Header
}

func getSearchResults(url string, token string) (GithubData, error) {
	log.Printf("sending request to: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GithubData{}, fmt.Errorf("error creating newrequest: %s", err.Error())
	}

	req.Header.Add("authorization", fmt.Sprintf("token %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GithubData{}, fmt.Errorf("error sending httprequest: %s", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GithubData{}, fmt.Errorf("error reading http response: %s", err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return GithubData{}, fmt.Errorf("github respondes with http status %d: %s", res.StatusCode, res.Status)
	}

	githubSearchData, err := UnmarshalGithubSearchData(body)
	if err != nil {
		return GithubData{}, fmt.Errorf("error unmarshalling http response: %s", err.Error())
	}

	return GithubData{
		Items:       githubSearchData.Items,
		HTTPHeaders: res.Header,
	}, nil
}

func getActivityContent(url string) (FlogoActivity, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FlogoActivity{}, fmt.Errorf("error creating newrequest: %s", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlogoActivity{}, fmt.Errorf("error sending httprequest: %s", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return FlogoActivity{}, fmt.Errorf("error reading http response: %s", err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return FlogoActivity{}, fmt.Errorf("github respondes with http status %d: %s", res.StatusCode, res.Status)
	}

	activity, err := UnmarshalFlogoActivity(body)
	if err != nil {
		return FlogoActivity{}, fmt.Errorf("error unmarshalling http response: %s", err.Error())
	}

	return activity, nil
}

func getRepoDetails(url string) (RepoDetails, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RepoDetails{}, fmt.Errorf("error creating newrequest: %s", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RepoDetails{}, fmt.Errorf("error sending httprequest: %s", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return RepoDetails{}, fmt.Errorf("error reading http response: %s", err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return RepoDetails{}, fmt.Errorf("github respondes with http status %d: %s", res.StatusCode, res.Status)
	}

	repoDetails, err := UnmarshalRepoDetails(body)
	if err != nil {
		return RepoDetails{}, fmt.Errorf("error unmarshalling http response: %s", err.Error())
	}

	return repoDetails, nil
}

func getMaxPages(h http.Header) int {
	p := 0

	links := linkheader.Parse(h.Get("Link"))

	for _, link := range links {
		if link.Rel == "last" {
			// Get the last page number
			// Link headers are formatted as https://api.github.com/search/code?q=flogo%2Bfilename%3Aactivity.json&page=33
			p, _ = strconv.Atoi(link.URL[strings.Index(link.URL, "page=")+5:])
		}

	}

	log.Printf("There are a total of %d pages", p)
	return p
}
