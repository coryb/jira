package jiracmd

import "github.com/go-jira/jira/jiracli"

func RegisterAllCommands() {
	// Issues
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "browse", Entry: CmdBrowseRegistry(), Aliases: []string{"b"}, Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "create", Entry: CmdCreateRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "edit", Entry: CmdEditRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "view", Entry: CmdViewRegistry(), Aliases: []string{"issue"}, Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "list", Entry: CmdListRegistry(), Aliases: []string{"ls"}, Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "comment", Entry: CmdCommentRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "subtask", Entry: CmdSubtaskRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "fields", Entry: CmdFieldsRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "issuetypes", Entry: CmdIssueTypesRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "createmeta", Entry: CmdCreateMetaRegistry(), Group: "Issues"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "editmeta", Entry: CmdEditMetaRegistry(), Group: "Issues"})

	// Transitions
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "transition", Entry: CmdTransitionRegistry(""), Aliases: []string{"trans"}, Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "transitions", Entry: CmdTransitionsRegistry("transitions"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "in-progress", Entry: CmdTransitionRegistry("Progress"), Aliases: []string{"prog", "progress"}, Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "start", Entry: CmdTransitionRegistry("start"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "stop", Entry: CmdTransitionRegistry("stop"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "todo", Entry: CmdTransitionRegistry("To Do"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "backlog", Entry: CmdTransitionRegistry("Backlog"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "done", Entry: CmdTransitionRegistry("Done"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "close", Entry: CmdTransitionRegistry("close"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "reopen", Entry: CmdTransitionRegistry("reopen"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "resolve", Entry: CmdTransitionRegistry("resolve"), Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "acknowledge", Entry: CmdTransitionRegistry("acknowledge"), Aliases: []string{"ack"}, Group: "Transitions"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "transmeta", Entry: CmdTransitionsRegistry("debug"), Group: "Transitions"})

	// Epics
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "epic create", Entry: CmdEpicCreateRegistry(), Group: "Epics"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "epic add", Entry: CmdEpicAddRegistry(), Group: "Epics"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "epic list", Entry: CmdEpicListRegistry(), Aliases: []string{"ls"}, Group: "Epics"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "epic remove", Entry: CmdEpicRemoveRegistry(), Aliases: []string{"rm"}, Group: "Epics"})

	// Attachments
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "attach create", Entry: CmdAttachCreateRegistry(), Group: "Attachments"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "attach get", Entry: CmdAttachGetRegistry(), Group: "Attachments"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "attach list", Entry: CmdAttachListRegistry(), Aliases: []string{"ls"}, Group: "Attachments"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "attach remove", Entry: CmdAttachRemoveRegistry(), Aliases: []string{"rm"}, Group: "Attachments"})

	// Labels
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "labels add", Entry: CmdLabelsAddRegistry(), Group: "Labels"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "labels remove", Entry: CmdLabelsRemoveRegistry(), Aliases: []string{"rm"}, Group: "Labels"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "labels set", Entry: CmdLabelsSetRegistry(), Group: "Labels"})

	// Linking
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "issuelink", Entry: CmdIssueLinkRegistry(), Group: "Linking"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "issuelinktypes", Entry: CmdIssueLinkTypesRegistry(), Group: "Linking"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "block", Entry: CmdBlockRegistry(), Group: "Linking"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "dup", Entry: CmdDupRegistry(), Group: "Linking"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "rank", Entry: CmdRankRegistry(), Group: "Linking"})

	// User
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "assign", Entry: CmdAssignRegistry(), Aliases: []string{"give"}, Group: "User"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "unassign", Entry: CmdUnassignRegistry(), Group: "User"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "take", Entry: CmdTakeRegistry(), Group: "User"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "watch", Entry: CmdWatchRegistry(), Group: "User"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "vote", Entry: CmdVoteRegistry(), Group: "User"})

	// Worklog
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "worklog add", Entry: CmdWorklogAddRegistry(), Group: "Worklog"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "worklog list", Entry: CmdWorklogListRegistry(), Default: true, Group: "Worklog"})

	// Components
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "components", Entry: CmdComponentsRegistry(), Group: "Components"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "component add", Entry: CmdComponentAddRegistry(), Group: "Components"})

	// Templates
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "export-templates", Entry: CmdExportTemplatesRegistry(), Group: "Templates"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "unexport-templates", Entry: CmdUnexportTemplatesRegistry(), Group: "Templates"})

	// Server
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "login", Entry: CmdLoginRegistry(), Group: "Server"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "logout", Entry: CmdLogoutRegistry(), Group: "Server"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "session", Entry: CmdSessionRegistry(), Group: "Server"})
	jiracli.RegisterCommand(jiracli.CommandRegistry{Command: "request", Entry: CmdRequestRegistry(), Aliases: []string{"req"}, Group: "Server"})
}
