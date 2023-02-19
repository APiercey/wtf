package github

import (
	"context"
	"fmt"

	ghb "github.com/google/go-github/v32/github"
)

func loadPullRequests(github *ghb.Client, username string) ([]*ghb.PullRequest, error) {
	query := fmt.Sprintf("is:open is:pr author:%s archived:false", username)

	issues, err := searchIssues(github, query)

	if err != nil {
		return []*ghb.PullRequest{}, err
	}

	return issuesToPullRequests(issues), nil
}

func loadReviewRequests(github *ghb.Client, username string) ([]*ghb.PullRequest, error) {
	query := fmt.Sprintf("is:open is:pr review-requested:%s archived:false ", username)

	issues, err := searchIssues(github, query)

	if err != nil {
		return []*ghb.PullRequest{}, err
	}

	return issuesToPullRequests(issues), nil
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

func issuesToPullRequests(issues []*ghb.Issue) []*ghb.PullRequest {
	prs := []*ghb.PullRequest{}

	for _, issue := range issues {
		if issue.IsPullRequest() {
			pr := ghb.PullRequest{ID: issue.ID, Number: issue.Number, Title: issue.Title, HTMLURL: issue.HTMLURL, MergeableState: nil}
			prs = append(prs, &pr)
		}
	}

	return prs
}
