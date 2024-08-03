package assignedIssues

import (
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hghtwr/jiracli/customHelp"
	"github.com/hghtwr/jiracli/jiraApi"
	"github.com/hghtwr/jiracli/navigation"
	"github.com/hghtwr/jiracli/notifications"
	"github.com/hghtwr/jiracli/views/layout"
)

type AssignedIssuesModel struct {
	issueTable table.Model
	NavTo 		navigation.ScreenId
	Context 	navigation.Context
}


func CreateInitModel() AssignedIssuesModel{
	return AssignedIssuesModel{
		issueTable: createTable(),
		NavTo: navigation.AssignedIssueView,

	}
}

func (m AssignedIssuesModel) Init() tea.Cmd {
	return nil
}
func (m AssignedIssuesModel) View() string {

	help := customHelp.CreateDefaultHelp()
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.issueTable.View(),
		help.View(CreateAssignedIssuesKeyMap(&m.issueTable.KeyMap)),
	)
}

func (m AssignedIssuesModel) Update(msg tea.Msg) (navigation.ScreenModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
			case "d": //View details
				m.NavTo = navigation.IssueDetailView
				m.Context = navigation.Context{IssueId: m.issueTable.SelectedRow()[0]}
				return m, navigation.RefreshView() // We send a refreshcmd to the main model to update the view after navigation.
				// Otherwise, if all other cmds are processed, the view will only update on the next update cycle.
			}
	}

	switch msg := msg.(type) {
		/*case tea.WindowSizeMsg:

			m.issueTable.SetWidth(msg.Width)
			m.issueTable.SetHeight(msg.Height - 20)
			return m, nil*/
		case jiraApi.MyIssuesResponse:
			m.SetTableRows(msg.Issues)
			cmds = append(cmds, notifications.CreateNotificationMsg("Fetched issues", notifications.Success, notifications.Tray))
			return m, tea.Batch(cmds...)
		case notifications.NotificationCmd:
		  return m, nil

		default:

			m.issueTable, cmd = m.issueTable.Update(msg)
			cmds = append(cmds, cmd)
			cmds = append(cmds, jiraApi.FetchAssignedIssues())
			cmds = append(cmds, notifications.CreateNotificationMsg("Fetching assigned issues", notifications.Info, notifications.Bar))
			return m, tea.Batch(cmds...)

		}
}

func (m AssignedIssuesModel) GetContext() navigation.Context {
	return m.Context
}

func (m AssignedIssuesModel) GetNavTo() navigation.ScreenId {
	return m.NavTo
}
func (m AssignedIssuesModel) SetNavTo(navTo navigation.ScreenId) navigation.ScreenModel {
	m.NavTo = navTo
	return m
}

func (m AssignedIssuesModel) SetContext(context navigation.Context) navigation.ScreenModel {
	m.Context = context
	return m
}

func (m *AssignedIssuesModel) SetTableRows(issues []jira.Issue) {

	var menuRows []table.Row
	for _, issue := range issues {
		var priority string
		if issue.Fields.Priority.Name == "Low" || issue.Fields.Priority.Name == "Lowest" {
			priority = "▼ " + issue.Fields.Priority.Name
		} else if issue.Fields.Priority.Name == "High" || issue.Fields.Priority.Name == "Highest" {
			priority = "▲ " + issue.Fields.Priority.Name
		} else {
			priority = "► " + issue.Fields.Priority.Name
		}
		menuRows = append(menuRows, table.Row{
			issue.Key,
			issue.Fields.Summary,
			issue.Fields.Status.Name,
			priority,
			issue.Fields.Assignee.DisplayName,
		})
	}
	m.issueTable.SetRows(menuRows)
}

func createTable() table.Model {

	return table.New(
		table.WithFocused(true),
		table.WithColumns([]table.Column{
			{
				Title: "Id",
				Width: layout.GetWidthFraction(2),
			},
			{
				Title: "Title",
				Width: layout.GetWidthFraction(14),
			},
			{
				Title: "Status",
				Width: layout.GetWidthFraction(2),
			},
			{
				Title: "Priority",
				Width: layout.GetWidthFraction(2),
			},
			{
				Title: "Assignee",
				Width: layout.GetWidthFraction(4),
			},
		}),
	)}

func CreateAssignedIssuesKeyMap(tableKeyMap *table.KeyMap ) *customHelp.DefaultKeyMap {
	return customHelp.CreateDefaultKeyMap([][]key.Binding{
		{
			tableKeyMap.LineDown,
			tableKeyMap.LineUp,
			tableKeyMap.GotoBottom,
			tableKeyMap.GotoTop,
			tableKeyMap.PageDown,
			tableKeyMap.PageUp,
			tableKeyMap.HalfPageDown,
			tableKeyMap.HalfPageUp,
		},{
			key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "View details"),
			),
		},
		{
			key.NewBinding(
				key.WithKeys("m"),
				key.WithHelp("m", "Comment issue"),
			),
			key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "Change status"),
			),
			key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp("p", "Change priority"),
			),
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "Change assignee"),
			),
		},
})
}
