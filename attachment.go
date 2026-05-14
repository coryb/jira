package jira

import (
	"context"
	"io"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

func (j *Jira) GetAttachment(id string) (*models.IssueAttachmentMetadataScheme, error) {
	return GetAttachment(j.UA, j.Endpoint, id)
}

func GetAttachment(ua HttpClient, endpoint string, id string) (*models.IssueAttachmentMetadataScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, resp, err := client.Issue.Attachment.Metadata(context.Background(), id)
	return result, atlassianResponseError(resp, err)
}

func (j *Jira) DownloadAttachment(id string) (io.Reader, error) {
	return DownloadAttachment(j.UA, j.Endpoint, id)
}

func DownloadAttachment(ua HttpClient, endpoint string, id string) (io.Reader, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	resp, err := client.Issue.Attachment.Download(context.Background(), id, true)
	if err != nil {
		return nil, err
	}
	return &resp.Bytes, nil
}

func (j *Jira) RemoveAttachment(id string) error {
	return RemoveAttachment(j.UA, j.Endpoint, id)
}

func RemoveAttachment(ua HttpClient, endpoint string, id string) error {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Issue.Attachment.Delete(context.Background(), id)
	return atlassianResponseError(resp, err)
}
