package jiracmd

import (
	"strings"

	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
)

// ensureServerInfo fetches serverInfo if not yet cached on globals, populating
// JiraDeploymentType and BaseURL. It is a no-op if JiraDeploymentType is already set.
func ensureServerInfo(o *oreo.Client, globals *jiracli.GlobalOptions) error {
	if globals.JiraDeploymentType.Value != "" {
		return nil
	}
	serverInfo, err := jira.GetServerInfo(o, globals.Endpoint.Value)
	if err != nil {
		return err
	}
	globals.JiraDeploymentType.Value = strings.ToLower(serverInfo.DeploymentType)
	if serverInfo.BaseURL != "" && globals.BaseURL.Value == "" {
		globals.BaseURL.Value = serverInfo.BaseURL
	}
	return nil
}
