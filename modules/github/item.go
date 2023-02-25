package github

import (
	ghb "github.com/google/go-github/v32/github"
)

const (
	AuthoredPullRequest int = 0
	ReviewRequested         = 1
)

type Item struct {
	PullRequest *ghb.PullRequest
	ItemType int
	ID int
}

func anyReviewRequests(items []Item) bool {
	for _, item := range items {
		if item.ItemType == ReviewRequested {
			return true
		}
	}

	return false
}

func anyAuthoriedRequests(items []Item) bool {
	for _, item := range items {
		if item.ItemType == AuthoredPullRequest {
			return true
		}
	}

	return false
}
