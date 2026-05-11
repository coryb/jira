package jira

import (
	"strings"

	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// ListOfAttachment is a named slice that implements sort.Interface.
type ListOfAttachment []*models.IssueAttachmentScheme

func (l *ListOfAttachment) Len() int      { return len(*l) }
func (l *ListOfAttachment) Less(i, j int) bool { return (*l)[i].ID < (*l)[j].ID }
func (l *ListOfAttachment) Swap(i, j int)      { (*l)[i], (*l)[j] = (*l)[j], (*l)[i] }

type AllowedValues []interface{}

type Operations []string

type JSONType struct {
	Custom   string `json:"custom,omitempty" yaml:"custom,omitempty"`
	CustomID int    `json:"customId,omitempty" yaml:"customId,omitempty"`
	Items    string `json:"items,omitempty" yaml:"items,omitempty"`
	System   string `json:"system,omitempty" yaml:"system,omitempty"`
	Type     string `json:"type,omitempty" yaml:"type,omitempty"`
}

type FieldMeta struct {
	AllowedValues   AllowedValues `json:"allowedValues,omitempty" yaml:"allowedValues,omitempty"`
	AutoCompleteURL string        `json:"autoCompleteUrl,omitempty" yaml:"autoCompleteUrl,omitempty"`
	DefaultValue    interface{}   `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	HasDefaultValue bool          `json:"hasDefaultValue,omitempty" yaml:"hasDefaultValue,omitempty"`
	Key             string        `json:"key,omitempty" yaml:"key,omitempty"`
	Name            string        `json:"name,omitempty" yaml:"name,omitempty"`
	Operations      Operations    `json:"operations,omitempty" yaml:"operations,omitempty"`
	Required        bool          `json:"required,omitempty" yaml:"required,omitempty"`
	Schema          *JSONType     `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type FieldMetaMap map[string]*FieldMeta

type FieldOperation map[string]interface{}

type FieldOperations []FieldOperation

type FieldOperationsMap map[string]FieldOperations

type EntityProperty struct {
	Key   string      `json:"key,omitempty" yaml:"key,omitempty"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

type Properties []*EntityProperty

type Visibility struct {
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// Issue embeds go-atlassian's typed IssueSchemeV2 for standard field access.
// The outer Fields map shadows IssueSchemeV2.Fields in JSON so templates receive
// all fields including custom customfield_* entries as a flat map.
type Issue struct {
	*models.IssueSchemeV2
	Fields map[string]interface{} `json:"fields,omitempty" yaml:"fields,omitempty"`
}

type Issues []*Issue

type IssueUpdate struct {
	Fields          map[string]interface{} `json:"fields,omitempty" yaml:"fields,omitempty"`
	HistoryMetadata map[string]interface{} `json:"historyMetadata,omitempty" yaml:"historyMetadata,omitempty"`
	Properties      Properties             `json:"properties,omitempty" yaml:"properties,omitempty"`
	Transition      *Transition            `json:"transition,omitempty" yaml:"transition,omitempty"`
	Update          FieldOperationsMap     `json:"update,omitempty" yaml:"update,omitempty"`
}

type EditMeta struct {
	Fields FieldMetaMap `json:"fields,omitempty" yaml:"fields,omitempty"`
}


type Transition struct {
	Expand    string       `json:"expand,omitempty" yaml:"expand,omitempty"`
	Fields    FieldMetaMap `json:"fields,omitempty" yaml:"fields,omitempty"`
	HasScreen bool         `json:"hasScreen,omitempty" yaml:"hasScreen,omitempty"`
	ID        string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name      string       `json:"name,omitempty" yaml:"name,omitempty"`
	To        *models.StatusScheme `json:"to,omitempty" yaml:"to,omitempty"`
}

type Transitions []*Transition

func (t Transitions) Find(name string) *Transition {
	name = strings.ToLower(name)
	matches := Transitions{}
	for _, trans := range t {
		if strings.ToLower(trans.Name) == name {
			return trans
		}
		if strings.Contains(strings.ToLower(trans.Name), name) {
			matches = append(matches, trans)
		}
	}
	if len(matches) > 0 {
		return matches[0]
	}
	return nil
}

type TransitionsMeta struct {
	Expand      string      `json:"expand,omitempty" yaml:"expand,omitempty"`
	Transitions Transitions `json:"transitions,omitempty" yaml:"transitions,omitempty"`
}

type IssueType struct {
	AvatarID    int          `json:"avatarId,omitempty" yaml:"avatarId,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Fields      FieldMetaMap `json:"fields,omitempty" yaml:"fields,omitempty"`
	IconURL     string       `json:"iconUrl,omitempty" yaml:"iconUrl,omitempty"`
	ID          string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string       `json:"name,omitempty" yaml:"name,omitempty"`
	Self        string       `json:"self,omitempty" yaml:"self,omitempty"`
	Subtask     bool         `json:"subtask,omitempty" yaml:"subtask,omitempty"`
}

type IssueTypes []*IssueType

type CreateMetaProject struct {
	AvatarUrls map[string]string `json:"avatarUrls,omitempty" yaml:"avatarUrls,omitempty"`
	Expand     string            `json:"expand,omitempty" yaml:"expand,omitempty"`
	ID         string            `json:"id,omitempty" yaml:"id,omitempty"`
	IssueTypes IssueTypes        `json:"issuetypes,omitempty" yaml:"issuetypes,omitempty"`
	Key        string            `json:"key,omitempty" yaml:"key,omitempty"`
	Name       string            `json:"name,omitempty" yaml:"name,omitempty"`
	Self       string            `json:"self,omitempty" yaml:"self,omitempty"`
}

type Projects []*CreateMetaProject

type CreateMeta struct {
	Expand   string   `json:"expand,omitempty" yaml:"expand,omitempty"`
	Projects Projects `json:"projects,omitempty" yaml:"projects,omitempty"`
}

type Fields []string

type IssueRef struct {
	Fields *Fields `json:"fields,omitempty" yaml:"fields,omitempty"`
	ID     string  `json:"id,omitempty" yaml:"id,omitempty"`
	Key    string  `json:"key,omitempty" yaml:"key,omitempty"`
}

type LinkIssueRequest struct {
	Comment      *models.IssueCommentSchemeV2 `json:"comment,omitempty" yaml:"comment,omitempty"`
	InwardIssue  *IssueRef                    `json:"inwardIssue,omitempty" yaml:"inwardIssue,omitempty"`
	OutwardIssue *IssueRef                    `json:"outwardIssue,omitempty" yaml:"outwardIssue,omitempty"`
	Type         *models.LinkTypeScheme       `json:"type,omitempty" yaml:"type,omitempty"`
}

type RankRequest struct {
	Issues          []string `json:"issues,omitempty" yaml:"issues,omitempty"`
	RankBeforeIssue string   `json:"rankBeforeIssue,omitempty" yaml:"rankBeforeIssue,omitempty"`
	RankAfterIssue  string   `json:"rankAfterIssue,omitempty" yaml:"rankAfterIssue,omitempty"`
}

type SearchRequest struct {
	Fields     Fields `json:"fields,omitempty" yaml:"fields,omitempty"`
	JQL        string `json:"jql,omitempty" yaml:"jql,omitempty"`
	MaxResults int    `json:"maxResults,omitempty" yaml:"maxResults,omitempty"`
	StartAt    int    `json:"startAt,omitempty" yaml:"startAt,omitempty"`
}

type SearchResults struct {
	Expand     string            `json:"expand,omitempty" yaml:"expand,omitempty"`
	Issues     Issues            `json:"issues,omitempty" yaml:"issues,omitempty"`
	MaxResults int               `json:"maxResults,omitempty" yaml:"maxResults,omitempty"`
	Names      map[string]string `json:"names,omitempty" yaml:"names,omitempty"`
	StartAt    int               `json:"startAt,omitempty" yaml:"startAt,omitempty"`
	Total      int               `json:"total,omitempty" yaml:"total,omitempty"`
}
