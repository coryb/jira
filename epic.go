package jira

import (
	"context"

	jagile "github.com/ctreminiom/go-atlassian/v2/jira/agile"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

func newAgileClient(ua HttpClient, endpoint string) (*jagile.Client, error) {
	return jagile.New(ua, endpoint)
}

// https://docs.atlassian.com/jira-software/REST/latest/#agile/1.0/epic-getIssuesForEpic
func (j *Jira) EpicSearch(epic string, sp SearchProvider) (*models.BoardIssuePageScheme, error) {
	return EpicSearch(j.UA, j.Endpoint, epic, sp)
}

func EpicSearch(ua HttpClient, endpoint string, epic string, sp SearchProvider) (*models.BoardIssuePageScheme, error) {
	req := sp.ProvideSearchRequest()

	client, err := newAgileClient(ua, endpoint)
	if err != nil {
		return nil, err
	}

	opts := &models.IssueOptionScheme{
		JQL:    req.JQL,
		Fields: []string(req.Fields),
	}

	result, resp, err := client.Epic.Issues(context.Background(), epic, opts, req.StartAt, req.MaxResults)
	return result, atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira-software/REST/latest/#agile/1.0/epic-moveIssuesToEpic
func (j *Jira) EpicAddIssues(epic string, issues []string) error {
	return EpicAddIssues(j.UA, j.Endpoint, epic, issues)
}

func EpicAddIssues(ua HttpClient, endpoint string, epic string, issues []string) error {
	client, err := newAgileClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Epic.Move(context.Background(), epic, issues)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira-software/REST/latest/#agile/1.0/epic-removeIssuesFromEpic
func (j *Jira) EpicRemoveIssues(issues []string) error {
	return EpicRemoveIssues(j.UA, j.Endpoint, issues)
}

func EpicRemoveIssues(ua HttpClient, endpoint string, issues []string) error {
	client, err := newAgileClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Epic.Move(context.Background(), "none", issues)
	return atlassianResponseError(resp, err)
}
