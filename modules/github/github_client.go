package github

import (
	"context"
	"net/http"

	ghb "github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// client defines a new GitHub client structure
type GithubClient struct {
	apiKey    string
	baseURL   string
	uploadURL string
}

func NewGithubClient(apiKey, baseURL, uploadURL string) *GithubClient {
	return &GithubClient{
		apiKey:    apiKey,
		baseURL:   baseURL,
		uploadURL: uploadURL,
	}
}

/* -------------------- Unexported Functions -------------------- */

func (client *GithubClient) githubClient() (*ghb.Client, error) {
	oauthClient := client.oauthClient()

	if client.isGitHubEnterprise() {
		return ghb.NewEnterpriseClient(client.baseURL, client.uploadURL, oauthClient)
	}

	return ghb.NewClient(oauthClient), nil
}

func (client *GithubClient) isGitHubEnterprise() bool {
	if len(client.baseURL) > 0 {
		if client.uploadURL == "" {
			client.uploadURL = client.baseURL
		}
		return true
	}
	return false
}

func (client *GithubClient) oauthClient() *http.Client {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: client.apiKey},
	)

	return oauth2.NewClient(context.Background(), tokenService)
}
