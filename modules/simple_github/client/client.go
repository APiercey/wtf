package client

import (
	"context"
	"fmt"

	ghb "github.com/google/go-github/v32/github"
)

type PullRequest ghb.PullRequest

// client defines a new GitHub client structure
type GithubClient struct {
	apiKey    string
	baseURL   string
	uploadURL string
	github    *ghb.Client
}

// Returns a GithubClient or fails
func NewGithubClient(apiKey, baseURL, uploadURL string) GithubClient {
	if isGitHubEnterprise(baseURL) {
		uploadURL = baseURL
	}

	github, err := buildRealClient(apiKey, baseURL, uploadURL)

	if err != nil {
		panic(err.Error())
	}

	return GithubClient{
		apiKey:    apiKey,
		baseURL:   baseURL,
		uploadURL: uploadURL,
		github: github,
	}
}

/* -------------------- Exported Functions -------------------- */

func (client GithubClient) LoadPullRequests(username string) ([]*ghb.PullRequest, error) {
	query := fmt.Sprintf("is:open is:pr author:%s archived:false", username)

	issues, err := client.searchIssues(query)

	if err != nil {
		return []*ghb.PullRequest{}, err
	}

	return issuesToPullRequests(&client, issues), nil
}

func (client GithubClient) LoadReviewRequests(username string) ([]*ghb.PullRequest, error) {
	query := fmt.Sprintf("is:open is:pr review-requested:%s archived:false ", username)

	issues, err := client.searchIssues(query)

	if err != nil {
		return []*ghb.PullRequest{}, err
	}

	return issuesToPullRequests(&client, issues), nil
}

/* -------------------- Unexported Functions -------------------- */


func (client GithubClient) searchIssues(query string) ([]*ghb.Issue, error) {
	opts := &ghb.SearchOptions{}
	opts.ListOptions.PerPage = 100

	results, _, err := client.github.Search.Issues(context.Background(), query, opts)

	if err != nil {
		return []*ghb.Issue{}, err
	}

	return results.Issues, nil
}
