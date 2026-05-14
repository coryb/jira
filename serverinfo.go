package jira

import (
	"context"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

func GetServerInfo(ua HttpClient, endpoint string) (*models.ServerInformationScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, resp, err := client.Server.Info(context.Background())
	return result, atlassianResponseError(resp, err)
}
