package jira

import (
	"context"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/component-createComponent
func (j *Jira) CreateComponent(payload *models.ComponentPayloadScheme) (*models.ComponentScheme, error) {
	return CreateComponent(j.UA, j.Endpoint, payload)
}

func CreateComponent(ua HttpClient, endpoint string, payload *models.ComponentPayloadScheme) (*models.ComponentScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, resp, err := client.Project.Component.Create(context.Background(), payload)
	return result, atlassianResponseError(resp, err)
}
