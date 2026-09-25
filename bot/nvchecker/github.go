package nvchecker

import (
	"context"

	"github.com/google/go-github/v92/github"
)

var ghClient *github.Client = nil

func InitGithubClient(pat string) error {
	client, err := github.NewClient(github.WithAuthToken(pat))
	if err != nil {
		return err
	}

	ghClient = client

	return nil
}

// ListFoxOrgs for pat valid check
func ListFoxOrgs() error {
	_, _, err := ghClient.Organizations.List(context.Background(), "weilinfox", nil)

	return err
}
