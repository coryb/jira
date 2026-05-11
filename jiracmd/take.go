package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdTakeRegistry() *jiracli.CommandRegistryEntry {
	opts := AssignOptions{}

	return &jiracli.CommandRegistryEntry{
		"Assign issue to yourself",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			if opts.Assignee == "" {
				if err := ensureServerInfo(o, globals); err != nil {
					return err
				}
				me, err := jira.GetCurrentUser(o, globals.Endpoint.Value)
				if err != nil {
					return err
				}
				if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
					opts.Assignee = me.AccountID
				} else {
					opts.Assignee = me.Name
				}
			}
			return CmdAssign(o, globals, &opts)
		},
	}
}

func CmdTakeUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	return nil
}
