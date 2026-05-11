package jira

import (
	"context"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/project-getProjectComponents
func (j *Jira) GetProjectComponents(project string) ([]*models.ComponentScheme, error) {
	return GetProjectComponents(j.UA, j.Endpoint, project)
}

func GetProjectComponents(ua HttpClient, endpoint string, project string) ([]*models.ComponentScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, _, err := client.Project.Component.Gets(context.Background(), project)
	return result, err
}

// https://developer.atlassian.com/cloud/jira/platform/rest/v2#api-api-2-project-projectIdOrKey-versions-get
func (j *Jira) GetProjectVersions(project string) ([]*models.VersionScheme, error) {
	return GetProjectVersions(j.UA, j.Endpoint, project)
}

func GetProjectVersions(ua HttpClient, endpoint string, project string) ([]*models.VersionScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, _, err := client.Project.Version.Gets(context.Background(), project)
	return result, err
}
