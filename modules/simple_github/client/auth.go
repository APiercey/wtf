package client

import (
	"context"
	"golang.org/x/oauth2"

	ghb "github.com/google/go-github/v32/github"
)

func buildRealClient(apiKey, baseURL, uploadURL string) (*ghb.Client, error) {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)

	oauthClient := oauth2.NewClient(context.Background(), tokenService)

	if isGitHubEnterprise(baseURL) {
		return ghb.NewEnterpriseClient(baseURL, uploadURL, oauthClient)
	}

	return ghb.NewClient(oauthClient), nil
}
