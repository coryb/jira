package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

type IssueQueryProvider interface {
	ProvideIssueQueryString() string
	ProvideFields() []string
	ProvideExpand() []string
}

type IssueOptions struct {
	Fields        []string `json:"fields,omitempty" yaml:"fields,omitempty"`
	Expand        []string `json:"expand,omitempty" yaml:"expand,omitempty"`
	Properties    []string `json:"properties,omitempty" yaml:"properties,omitempty"`
	FieldsByKeys  bool     `json:"fieldsByKeys,omitempty" yaml:"fieldsByKeys,omitempty"`
	UpdateHistory bool     `json:"updateHistory,omitempty" yaml:"updateHistory,omitempty"`
}

func (o *IssueOptions) ProvideFields() []string { return o.Fields }
func (o *IssueOptions) ProvideExpand() []string { return o.Expand }

func (o *IssueOptions) ProvideIssueQueryString() string {
	params := []string{}
	if len(o.Fields) > 0 {
		params = append(params, "fields="+strings.Join(o.Fields, ","))
	}
	if len(o.Expand) > 0 {
		params = append(params, "expand="+strings.Join(o.Expand, ","))
	}
	if len(o.Properties) > 0 {
		params = append(params, "properties="+strings.Join(o.Properties, ","))
	}
	if o.FieldsByKeys {
		params = append(params, "fieldsByKeys=true")
	}
	if o.UpdateHistory {
		params = append(params, "updateHistory=true")
	}
	if len(params) > 0 {
		return "?" + strings.Join(params, "&")
	}
	return ""
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getIssue
func (j *Jira) GetIssue(issue string, iqg IssueQueryProvider) (*Issue, error) {
	return GetIssue(j.UA, j.Endpoint, issue, iqg)
}

func GetIssue(ua HttpClient, endpoint string, issue string, iqg IssueQueryProvider) (*Issue, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	var fields, expand []string
	if iqg != nil {
		fields = iqg.ProvideFields()
		expand = iqg.ProvideExpand()
	}
	typed, resp, err := client.Issue.Get(context.Background(), issue, fields, expand)
	if err != nil {
		return nil, err
	}
	result := &Issue{IssueSchemeV2: typed}
	var raw struct {
		Fields map[string]interface{} `json:"fields"`
	}
	if jsonErr := json.Unmarshal(resp.Bytes.Bytes(), &raw); jsonErr == nil {
		result.Fields = raw.Fields
	}
	return result, nil
}

func (j *Jira) GetIssueWorklog(issue string) (*[]*models.IssueWorklogRichTextScheme, error) {
	return GetIssueWorklog(j.UA, j.Endpoint, issue)
}

func GetIssueWorklog(ua HttpClient, endpoint string, issue string) (*[]*models.IssueWorklogRichTextScheme, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	page, resp, err := client.Issue.Worklog.Issue(context.Background(), issue, 0, 0, 0, nil)
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	return &page.Worklogs, nil
}

func (j *Jira) GetIssueComment(issue string) ([]*models.IssueCommentSchemeV2, error) {
	return GetIssueComment(j.UA, j.Endpoint, issue)
}

// https://docs.atlassian.com/software/jira/docs/api/REST/7.12.0/#api/2/issue-getComments
func GetIssueComment(ua HttpClient, endpoint string, issue string) ([]*models.IssueCommentSchemeV2, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	startAt := 0
	total := 1
	maxResults := 100
	var comments []*models.IssueCommentSchemeV2
	for startAt < total {
		page, resp, err := client.Issue.Comment.Gets(context.Background(), issue, "", nil, startAt, maxResults)
		if err != nil {
			return nil, atlassianResponseError(resp, err)
		}
		total = page.Total
		comments = append(comments, page.Comments...)
		startAt += maxResults
	}
	return comments, nil
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/worklog-addWorklog
func (j *Jira) AddIssueWorklog(issue string, req *models.IssueWorklogRichTextScheme) (*models.IssueWorklogRichTextScheme, error) {
	return AddIssueWorklog(j.UA, j.Endpoint, issue, req)
}

func AddIssueWorklog(ua HttpClient, endpoint string, issue string, req *models.IssueWorklogRichTextScheme) (*models.IssueWorklogRichTextScheme, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	payload := &models.WorklogRichTextPayloadScheme{
		Started:          req.Started,
		TimeSpent:        req.TimeSpent,
		TimeSpentSeconds: req.TimeSpentSeconds,
		Visibility:       req.Visibility,
	}
	if req.Comment != "" {
		payload.Comment = &models.CommentPayloadSchemeV2{Body: req.Comment}
	}
	result, resp, err := client.Issue.Worklog.Add(context.Background(), issue, payload, nil)
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	return result, nil
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getEditIssueMeta
func (j *Jira) GetIssueEditMeta(issue string) (*EditMeta, error) {
	return GetIssueEditMeta(j.UA, j.Endpoint, issue)
}

func GetIssueEditMeta(ua HttpClient, endpoint string, issue string) (*EditMeta, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	_, resp, err := client.Issue.Metadata.Get(context.Background(), issue, false, false)
	if err != nil {
		return nil, err
	}
	results := &EditMeta{}
	return results, json.Unmarshal(resp.Bytes.Bytes(), results)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-editIssue
func (j *Jira) EditIssue(issue string, req *IssueUpdate) error {
	return EditIssue(j.UA, j.Endpoint, issue, req)
}

func EditIssue(ua HttpClient, endpoint string, issue string, req *IssueUpdate) error {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(encoded, &bodyMap); err != nil {
		return err
	}
	cf := &models.CustomFields{Fields: []map[string]interface{}{bodyMap}}
	resp, err := client.Issue.Update(context.Background(), issue, true, &models.IssueSchemeV2{}, cf, nil)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-createIssue
func (j *Jira) CreateIssue(req *IssueUpdate) (*models.IssueResponseScheme, error) {
	return CreateIssue(j.UA, j.Endpoint, req)
}

func CreateIssue(ua HttpClient, endpoint string, req *IssueUpdate) (*models.IssueResponseScheme, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(encoded, &bodyMap); err != nil {
		return nil, err
	}
	cf := &models.CustomFields{Fields: []map[string]interface{}{bodyMap}}
	result, resp, err := client.Issue.Create(context.Background(), &models.IssueSchemeV2{}, cf)
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	return result, nil
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getCreateIssueMeta
func (j *Jira) GetIssueCreateMetaProject(projectKey string) (*CreateMetaProject, error) {
	return GetIssueCreateMetaProject(j.UA, j.Endpoint, projectKey)
}

func GetIssueCreateMetaProject(ua HttpClient, endpoint string, projectKey string) (*CreateMetaProject, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	opts := &models.IssueMetadataCreateOptions{
		ProjectKeys: []string{projectKey},
		Expand:      "projects.issuetypes.fields",
	}
	_, resp, err := client.Issue.Metadata.Create(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	results := &CreateMeta{}
	if err := json.Unmarshal(resp.Bytes.Bytes(), results); err != nil {
		return nil, err
	}
	for _, project := range results.Projects {
		if project.Key == projectKey {
			return project, nil
		}
	}
	return nil, fmt.Errorf("project %s not found", projectKey)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getCreateIssueMeta
func (j *Jira) GetIssueCreateMetaIssueType(projectKey, issueTypeName string) (*IssueType, error) {
	return GetIssueCreateMetaIssueType(j.UA, j.Endpoint, projectKey, issueTypeName)
}

func GetIssueCreateMetaIssueType(ua HttpClient, endpoint string, projectKey, issueTypeName string) (*IssueType, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	opts := &models.IssueMetadataCreateOptions{
		ProjectKeys:    []string{projectKey},
		IssueTypeNames: []string{issueTypeName},
		Expand:         "projects.issuetypes.fields",
	}
	_, resp, err := client.Issue.Metadata.Create(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	results := &CreateMeta{}
	if err := json.Unmarshal(resp.Bytes.Bytes(), results); err != nil {
		return nil, err
	}
	for _, project := range results.Projects {
		if project.Key != projectKey {
			continue
		}
		for _, issueType := range project.IssueTypes {
			if issueType.Name == issueTypeName {
				return issueType, nil
			}
		}
	}
	return nil, fmt.Errorf("project %s and IssueType %s not found", projectKey, issueTypeName)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issueLink-linkIssues
func (j *Jira) LinkIssues(req *LinkIssueRequest) error {
	return LinkIssues(j.UA, j.Endpoint, req)
}

func LinkIssues(ua HttpClient, endpoint string, req *LinkIssueRequest) error {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return err
	}
	payload := &models.LinkPayloadSchemeV2{Type: req.Type}
	if req.InwardIssue != nil {
		payload.InwardIssue = &models.LinkedIssueScheme{ID: req.InwardIssue.ID, Key: req.InwardIssue.Key}
	}
	if req.OutwardIssue != nil {
		payload.OutwardIssue = &models.LinkedIssueScheme{ID: req.OutwardIssue.ID, Key: req.OutwardIssue.Key}
	}
	if req.Comment != nil {
		payload.Comment = &models.CommentPayloadSchemeV2{Body: req.Comment.Body, Visibility: req.Comment.Visibility}
	}
	resp, err := client.Issue.Link.Create(context.Background(), payload)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getTransitions
func (j *Jira) GetIssueTransitions(issue string) (*TransitionsMeta, error) {
	return GetIssueTransitions(j.UA, j.Endpoint, issue)
}

func GetIssueTransitions(ua HttpClient, endpoint string, issue string) (*TransitionsMeta, error) {
	uri := URLJoin(endpoint, "rest/api/2/issue", issue, "transitions")
	uri += "?expand=transitions.fields"
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &TransitionsMeta{}
		return results, json.NewDecoder(resp.Body).Decode(results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-doTransition
func (j *Jira) TransitionIssue(issue string, req *IssueUpdate) error {
	return TransitionIssue(j.UA, j.Endpoint, issue, req)
}

func TransitionIssue(ua HttpClient, endpoint string, issue string, req *IssueUpdate) error {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return err
	}
	transitionID := ""
	if req.Transition != nil {
		transitionID = req.Transition.ID
	}
	var options *models.IssueMoveOptionsV2
	if req.Fields != nil || req.Update != nil {
		bodyMap := map[string]interface{}{}
		if req.Fields != nil {
			bodyMap["fields"] = req.Fields
		}
		if req.Update != nil {
			bodyMap["update"] = req.Update
		}
		options = &models.IssueMoveOptionsV2{
			Fields:       &models.IssueSchemeV2{},
			CustomFields: &models.CustomFields{Fields: []map[string]interface{}{bodyMap}},
		}
	}
	resp, err := client.Issue.Move(context.Background(), issue, transitionID, options)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issueLinkType-getIssueLinkTypes
func (j *Jira) GetIssueLinkTypes() ([]*models.LinkTypeScheme, error) {
	return GetIssueLinkTypes(j.UA, j.Endpoint)
}

func GetIssueLinkTypes(ua HttpClient, endpoint string) ([]*models.LinkTypeScheme, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, resp, err := client.Issue.Link.Type.Gets(context.Background())
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	return result.IssueLinkTypes, nil
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-addVote
func (j *Jira) IssueAddVote(issue string) error {
	return IssueAddVote(j.UA, j.Endpoint, issue)
}

func IssueAddVote(ua HttpClient, endpoint string, issue string) error {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Issue.Vote.Add(context.Background(), issue)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-removeVote
func (j *Jira) IssueRemoveVote(issue string) error {
	return IssueRemoveVote(j.UA, j.Endpoint, issue)
}

func IssueRemoveVote(ua HttpClient, endpoint string, issue string) error {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Issue.Vote.Delete(context.Background(), issue)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira-software/REST/cloud/#agile/1.0/issue-rankIssues
func (j *Jira) RankIssues(req *RankRequest) error {
	return RankIssues(j.UA, j.Endpoint, req)
}

func RankIssues(ua HttpClient, endpoint string, req *RankRequest) error {
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}
	uri := URLJoin(endpoint, "rest/agile/1.0/issue/rank")
	resp, err := ua.Put(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-addWatcher
func (j *Jira) IssueAddWatcher(issue, user string) error {
	return IssueAddWatcher(j.UA, j.Endpoint, issue, user)
}

func IssueAddWatcher(ua HttpClient, endpoint string, issue, user string) error {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Issue.Watcher.Add(context.Background(), issue, user)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-addWatcher
func (j *Jira) IssueRemoveWatcher(issue, user string) error {
	return IssueRemoveWatcher(j.UA, j.Endpoint, issue, user)
}

func IssueRemoveWatcher(ua HttpClient, endpoint string, issue, user string) error {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Issue.Watcher.Delete(context.Background(), issue, user)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/comment-addComment
func (j *Jira) IssueAddComment(issue string, req *models.IssueCommentSchemeV2) (*models.IssueCommentSchemeV2, error) {
	return IssueAddComment(j.UA, j.Endpoint, issue, req)
}

func IssueAddComment(ua HttpClient, endpoint string, issue string, req *models.IssueCommentSchemeV2) (*models.IssueCommentSchemeV2, error) {
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return nil, err
	}
	payload := &models.CommentPayloadSchemeV2{Body: req.Body, Visibility: req.Visibility}
	result, resp, err := client.Issue.Comment.Add(context.Background(), issue, payload, nil)
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	return result, nil
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-assign
func (j *Jira) IssueAssign(issue, name string) error {
	return IssueAssign(j.UA, j.Endpoint, issue, name)
}

func IssueAssign(ua HttpClient, endpoint string, issue, name string) error {
	if name == "" {
		return IssueAssignAccountID(ua, endpoint, issue, "")
	}
	users, err := UserSearch(ua, endpoint, &UserSearchOptions{Query: name, MaxResults: 10})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("user %q not found", name)
	}
	return IssueAssignAccountID(ua, endpoint, issue, users[0].AccountID)
}

func IssueAssignAccountID(ua HttpClient, endpoint string, issue, acctId string) error {
	if acctId == "" {
		// Unassign: go-atlassian enforces non-empty accountId, so send null directly.
		req := struct {
			AccountID *string `json:"accountId"`
		}{nil}
		encoded, err := json.Marshal(req)
		if err != nil {
			return err
		}
		uri := URLJoin(endpoint, "rest/api/2/issue", issue, "assignee")
		resp, err := ua.Put(uri, "application/json", bytes.NewBuffer(encoded))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 204 {
			return nil
		}
		return responseError(resp)
	}
	client, err := newAtlassianClientV2(ua, endpoint)
	if err != nil {
		return err
	}
	resp, err := client.Issue.Assign(context.Background(), issue, acctId)
	return atlassianResponseError(resp, err)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/attachments-addAttachment
func (j *Jira) IssueAttachFile(issue, filename string, contents io.Reader) (*ListOfAttachment, error) {
	return IssueAttachFile(j.UA, j.Endpoint, issue, filename, contents)
}

func IssueAttachFile(ua HttpClient, endpoint string, issue, filename string, contents io.Reader) (*ListOfAttachment, error) {
	client, err := newAtlassianClient(ua, endpoint)
	if err != nil {
		return nil, err
	}
	result, resp, err := client.Issue.Attachment.Add(context.Background(), issue, filename, contents)
	if err != nil {
		return nil, atlassianResponseError(resp, err)
	}
	out := ListOfAttachment(result)
	return &out, nil
}
