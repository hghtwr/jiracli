package customHelp

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type DefaultKeyMap struct {
	HistoryBack key.Binding
	Quit key.Binding
	NavigationBindings []key.Binding
	CustomBindings [][]key.Binding
}

func CreateDefaultHelp() *help.Model {
	userHelp := help.New()
	userHelp.ShowAll = true
	userHelp.FullSeparator = "   |   "
	userHelp.Styles = help.Styles{
		FullKey: lipgloss.NewStyle().Foreground(lipgloss.Color("230")),
		FullDesc: lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		FullSeparator: lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	}

	return &userHelp
}

// You can extend the default keymap with as many custom keybindings as you like, each array will create a new column in the help view.
func CreateDefaultKeyMap(customBindings [][]key.Binding) *DefaultKeyMap {
	NavigationBindings:= []key.Binding{
		key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "Settings"),
		),
		key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "Assigned Issues"),
		),
	}

	historyBackBinding := key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "Go back"),
	)
	quitBinding := key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	)
	return &DefaultKeyMap{
		HistoryBack: historyBackBinding,
		Quit: quitBinding,
		CustomBindings: customBindings,
		NavigationBindings: NavigationBindings,
	}
}



func (k DefaultKeyMap) ShortHelp() []key.Binding {
	bindings := []key.Binding{k.HistoryBack, k.Quit}
	return append(bindings, k.NavigationBindings...)
}


func (k DefaultKeyMap) FullHelp() [][]key.Binding {
	bindings := [][]key.Binding{}
	bindings = append(bindings, k.NavigationBindings)
	bindings = append(bindings, k.CustomBindings...)
	return append(bindings, []key.Binding{k.HistoryBack, k.Quit})
}

/*

// Refresh is not necessary anymore but we can reuse the custom keybinding for other options, e.g. filtering, sorting, moving issues to other states, etc.
type assignedIssuesKeyMap struct {
	CommentIssue key.Binding
	ChangeStatus key.Binding
	ChangePriority key.Binding
	ChangeAssignee key.Binding
	table.KeyMap
}

func (k assignedIssuesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CommentIssue, k.LineUp, k.LineDown }
}

func (k assignedIssuesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown, k.GotoBottom, k.GotoTop, k.PageDown, k.PageUp },
		{k.CommentIssue, k.ChangeStatus, k.ChangePriority, k.ChangeAssignee},
	}
}

func CreateIssueViewKeyMap() *assignedIssuesKeyMap {

	commentBinding := key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Comment issue"),
	)
	statusBinding := key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Change status"),
	)
	priorityBinding := key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Change priority"),
	)
	assigneeBinding := key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Change assignee"),
	)

	return &assignedIssuesKeyMap{
		CommentIssue: commentBinding,
		ChangeStatus: statusBinding,
		ChangePriority: priorityBinding,
		ChangeAssignee: assigneeBinding,
		//KeyMap: *tableKeyMap,
	}
}
*/