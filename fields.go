package jira

import (
	"context"

	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

func (j *Jira) GetFields() ([]*models.IssueFieldScheme, error) {
	return GetFields(j.UA, j.Endpoint)
}

func GetFields(ua HttpClient, endpoint string) ([]*models.IssueFieldScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, resp, err := client.Issue.Field.Gets(context.Background())
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	return result, nil
}
