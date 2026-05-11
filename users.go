package jira

import (
	"context"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

type UserSearchOptions struct {
	Query      string `yaml:"query,omitempty" json:"query,omitempty"`
	Username   string `yaml:"username,omitempty" json:"username,omitempty"`
	AccountID  string `yaml:"accountId,omitempty" json:"accountId,omitempty"`
	StartAt    int    `yaml:"startAt,omitempty" json:"startAt,omitempty"`
	MaxResults int    `yaml:"max-results,omitempty" json:"max-results,omitempty"`
	Property   string `yaml:"property,omitempty" json:"property,omitempty"`
}

// https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-user-search-get

func UserSearch(ua HttpClient, endpoint string, opts *UserSearchOptions) ([]*models.UserScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, _, err := client.User.Search.Do(context.Background(), opts.AccountID, opts.Query, opts.StartAt, opts.MaxResults)
	return result, err
}

func GetCurrentUser(ua HttpClient, endpoint string) (*models.UserScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, _, err := client.MySelf.Details(context.Background(), nil)
	return result, err
}
