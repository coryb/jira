package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EpicAddOptions struct {
	Issues  []string `yaml:"issues,omitempty"  json:"issues,omitempty"`
	Project string   `yaml:"project,omitempty" json:"project,omitempty"`
	Epic    string   `yaml:"epic,omitempty"    json:"epic,omitempty"`
}

func CmdEpicAddRegistry() *jiracli.CommandRegistryEntry {
	opts := EpicAddOptions{}

	return &jiracli.CommandRegistryEntry{
		"Add issues to Epic",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEpicAddUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Epic = jiracli.FormatIssue(opts.Epic, opts.Project)
			for i := range opts.Issues {
				opts.Issues[i] = jiracli.FormatIssue(opts.Issues[i], opts.Project)
			}
			return CmdEpicAdd(o, globals, &opts)
		},
	}
}

func CmdEpicAddUsage(cmd *kingpin.CmdClause, opts *EpicAddOptions) error {
	cmd.Arg("EPIC", "Epic Key or ID to add issues to").Required().StringVar(&opts.Epic)
	cmd.Arg("ISSUE", "Issues to add to epic").Required().StringsVar(&opts.Issues)
	return nil
}

func CmdEpicAdd(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EpicAddOptions) error {
	if err := jira.EpicAddIssues(o, globals.Endpoint.Value, opts.Epic, opts.Issues); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.Epic, globals.BrowseURL(opts.Epic))
		for _, issue := range opts.Issues {
			fmt.Printf("OK %s %s\n", issue, globals.BrowseURL(issue))
		}
	}

	return nil
}
