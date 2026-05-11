package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	models "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type SubtaskOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jira.IssueUpdate  `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string            `yaml:"project,omitempty" json:"project,omitempty"`
	IssueType             string            `yaml:"issuetype,omitempty" json:"issuetype,omitempty"`
	Overrides             map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Issue                 string            `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdSubtaskRegistry() *jiracli.CommandRegistryEntry {
	opts := SubtaskOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("subtask"),
		},
		Overrides: map[string]string{},
	}

	return &jiracli.CommandRegistryEntry{
		"Subtask issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdSubtaskUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			if opts.IssueType == "" {
				opts.IssueType = "Sub-task"
			}
			return CmdSubtask(o, globals, &opts)
		},
	}
}

func CmdSubtaskUsage(cmd *kingpin.CmdClause, opts *SubtaskOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("project", "project to subtask issue in").Short('p').StringVar(&opts.Project)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Arg("ISSUE", "Parent issue for subtask").StringVar(&opts.Issue)
	return nil
}

// CmdSubtask sends the subtask-metadata to the "subtask" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func CmdSubtask(o *oreo.Client, globals *jiracli.GlobalOptions, opts *SubtaskOptions) error {
	if err := ensureServerInfo(o, globals); err != nil {
		return err
	}

	type templateInput struct {
		Meta      *jira.IssueType `yaml:"meta" json:"meta"`
		Overrides map[string]string   `yaml:"overrides" json:"overrides"`
		Parent    *jira.Issue     `yaml:"parent" json:"parent"`
	}

	parent, err := jira.GetIssue(o, globals.Endpoint.Value, opts.Issue, nil)
	if err != nil {
		return err
	}

	if project, ok := parent.Fields["project"].(map[string]interface{}); ok {
		if key, ok := project["key"].(string); ok {
			opts.Project = key
		} else {
			return fmt.Errorf("Failed to find Project Key in parent issue")
		}
	} else {
		return fmt.Errorf("Failed to find Project field in parent issue")
	}

	createMeta, err := jira.GetIssueCreateMetaIssueType(o, globals.Endpoint.Value, opts.Project, opts.IssueType)
	if err != nil {
		return err
	}

	issueUpdate := jira.IssueUpdate{}
	input := templateInput{
		Meta:      createMeta,
		Overrides: opts.Overrides,
		Parent:    parent,
	}
	input.Overrides["project"] = opts.Project
	input.Overrides["issuetype"] = opts.IssueType
	input.Overrides["login"] = globals.Login.Value
	if me, err := jira.GetCurrentUser(o, globals.Endpoint.Value); err == nil && me.DisplayName != "" {
		input.Overrides["displayName"] = me.DisplayName
	}

	var issueResp *models.IssueResponseScheme
	err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
		if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
			err := fixGDPRUserFields(o, globals.Endpoint.Value, createMeta.Fields, issueUpdate.Fields)
			if err != nil {
				return err
			}
		}
		issueResp, err = jira.CreateIssue(o, globals.Endpoint.Value, &issueUpdate)
		return err
	})
	if err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", issueResp.Key, globals.BrowseURL(issueResp.Key))
	}

	if opts.Browse.Value {
		return CmdBrowse(globals, issueResp.Key)
	}
	return nil
}
