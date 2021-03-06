package github

import "encoding/json"

func UnmarshalGithubSearchData(data []byte) (GithubSearchData, error) {
	var r GithubSearchData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GithubSearchData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GithubSearchData struct {
	TotalCount        int64  `json:"total_count"`
	IncompleteResults bool   `json:"incomplete_results"`
	Items             []Item `json:"items"`
}

type Item struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	SHA        string     `json:"sha"`
	URL        string     `json:"url"`
	GitURL     string     `json:"git_url"`
	HTMLURL    string     `json:"html_url"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	ID               int64   `json:"id"`
	NodeID           string  `json:"node_id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	Private          bool    `json:"private"`
	Owner            Owner   `json:"owner"`
	HTMLURL          string  `json:"html_url"`
	Description      *string `json:"description"`
	Fork             bool    `json:"fork"`
	URL              string  `json:"url"`
	ForksURL         string  `json:"forks_url"`
	KeysURL          string  `json:"keys_url"`
	CollaboratorsURL string  `json:"collaborators_url"`
	TeamsURL         string  `json:"teams_url"`
	HooksURL         string  `json:"hooks_url"`
	IssueEventsURL   string  `json:"issue_events_url"`
	EventsURL        string  `json:"events_url"`
	AssigneesURL     string  `json:"assignees_url"`
	BranchesURL      string  `json:"branches_url"`
	TagsURL          string  `json:"tags_url"`
	BlobsURL         string  `json:"blobs_url"`
	GitTagsURL       string  `json:"git_tags_url"`
	GitRefsURL       string  `json:"git_refs_url"`
	TreesURL         string  `json:"trees_url"`
	StatusesURL      string  `json:"statuses_url"`
	LanguagesURL     string  `json:"languages_url"`
	StargazersURL    string  `json:"stargazers_url"`
	ContributorsURL  string  `json:"contributors_url"`
	SubscribersURL   string  `json:"subscribers_url"`
	SubscriptionURL  string  `json:"subscription_url"`
	CommitsURL       string  `json:"commits_url"`
	GitCommitsURL    string  `json:"git_commits_url"`
	CommentsURL      string  `json:"comments_url"`
	IssueCommentURL  string  `json:"issue_comment_url"`
	ContentsURL      string  `json:"contents_url"`
	CompareURL       string  `json:"compare_url"`
	MergesURL        string  `json:"merges_url"`
	ArchiveURL       string  `json:"archive_url"`
	DownloadsURL     string  `json:"downloads_url"`
	IssuesURL        string  `json:"issues_url"`
	PullsURL         string  `json:"pulls_url"`
	MilestonesURL    string  `json:"milestones_url"`
	NotificationsURL string  `json:"notifications_url"`
	LabelsURL        string  `json:"labels_url"`
	ReleasesURL      string  `json:"releases_url"`
	DeploymentsURL   string  `json:"deployments_url"`
}

type Owner struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

func UnmarshalRepoDetails(data []byte) (RepoDetails, error) {
	var r RepoDetails
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RepoDetails) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RepoDetails struct {
	ID               int64       `json:"id"`
	NodeID           string      `json:"node_id"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Private          bool        `json:"private"`
	Owner            Owner       `json:"owner"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
	PushedAt         string      `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         interface{} `json:"homepage"`
	Size             int64       `json:"size"`
	StargazersCount  int64       `json:"stargazers_count"`
	WatchersCount    int64       `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int64       `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int64       `json:"open_issues_count"`
	License          interface{} `json:"license"`
	Forks            int64       `json:"forks"`
	OpenIssues       int64       `json:"open_issues"`
	Watchers         int64       `json:"watchers"`
	DefaultBranch    string      `json:"default_branch"`
	TempCloneToken   interface{} `json:"temp_clone_token"`
	NetworkCount     int64       `json:"network_count"`
	SubscribersCount int64       `json:"subscribers_count"`
}
