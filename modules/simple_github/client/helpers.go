package client

import (
	"context"
	"strings"
	"strconv"
	"fmt"

	ghb "github.com/google/go-github/v32/github"
)

func issuesToPullRequests(client *GithubClient, issues []*ghb.Issue) []*ghb.PullRequest {
	prs := []*ghb.PullRequest{}

	for _, issue := range issues {

		if issue.IsPullRequest() {
			ownerName, repoName, prNumber := pullRequestInfo(issue)

			pr, _, err := client.github.PullRequests.Get(
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

func isGitHubEnterprise(baseURL string) bool {
	return len(baseURL) > 0
}
