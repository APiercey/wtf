package github

import (
	"context"
	"strings"
	"fmt"
	"strconv"

	ghb "github.com/google/go-github/v32/github"
)

func loadPullRequests(github *ghb.Client, username string) ([]*ghb.PullRequest, error) {
	query := fmt.Sprintf("is:open is:pr author:%s archived:false", username)

	issues, err := searchIssues(github, query)

	if err != nil {
		return []*ghb.PullRequest{}, err
	}

	return issuesToPullRequests(github, issues), nil
}

func loadReviewRequests(github *ghb.Client, username string) ([]*ghb.PullRequest, error) {
	query := fmt.Sprintf("is:open is:pr review-requested:%s archived:false ", username)

	issues, err := searchIssues(github, query)

	if err != nil {
		return []*ghb.PullRequest{}, err
	}

	return issuesToPullRequests(github, issues), nil
}

func searchIssues(client *ghb.Client, query string) ([]*ghb.Issue, error) {
	opts := &ghb.SearchOptions{}
	opts.ListOptions.PerPage = 100

	results, _, err := client.Search.Issues(context.Background(), query, opts)

	if err != nil {
		return []*ghb.Issue{}, err
	}

	return results.Issues, nil
}

func issuesToPullRequests(github *ghb.Client, issues []*ghb.Issue) []*ghb.PullRequest {
	prs := []*ghb.PullRequest{}

	for _, issue := range issues {

		if issue.IsPullRequest() {
			ownerName, repoName, prNumber := pullRequestInfo(issue)

			pr, _, err := github.PullRequests.Get(
				context.Background(),
				ownerName,
				repoName,
				prNumber,
			)

			if err == nil {
				prs = append(prs, pr)
			}
		}
	}

	return prs
}

func pullRequestInfo(issue *ghb.Issue) (ownerName string, repoName string, number int) {
	data := strings.Split(issue.PullRequestLinks.GetURL(), "repos")
	data = strings.Split(data[1], "/")

	pullRequestNumber, err := strconv.Atoi(data[4])

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	return data[1], data[2], pullRequestNumber
}
