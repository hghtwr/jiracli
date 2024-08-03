package changeAssignee

import (
	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hghtwr/jiracli/customHelp"
	"github.com/hghtwr/jiracli/jiraApi"
	"github.com/hghtwr/jiracli/navigation"
	"github.com/hghtwr/jiracli/notifications"
	"github.com/hghtwr/jiracli/views/layout"
)


type ChangeAssigneeModel struct {
	NavTo 		navigation.ScreenId
	Context 	navigation.Context
	issue     *jira.Issue
}

func CreateInitModel() ChangeAssigneeModel{
	return ChangeAssigneeModel{
		NavTo: navigation.ChangeAssigneeView,
	}
}

func (m ChangeAssigneeModel) Init() tea.Cmd {
	return nil
}

func (m ChangeAssigneeModel) View() string {
	headerStyle := lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.ThickBorder()).Width(layout.GetWidthFraction(24))
	headerContent := make([]string, 0)

	if m.issue != nil {
		headerContent = append(headerContent, layout.Style.HeaderStyle.Render(m.issue.Key))
		headerContent = append(headerContent, m.issue.Fields.Summary)
	}

	help := customHelp.CreateDefaultHelp()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		headerStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left, headerContent...)),
		help.View(customHelp.CreateDefaultKeyMap(nil)))
}

func (m ChangeAssigneeModel) Update(msg tea.Msg) (navigation.ScreenModel, tea.Cmd) {
	//var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case jiraApi.IssueDetailResponse:
		m.issue = msg.Issue
		//Now fetch tasks that have this task as parent (different from subtasks, e.g. for stories in epics)
		return m, tea.Batch(cmds...)
	default:
		cmd := jiraApi.FetchIssueDetails(m.Context.IssueId)
		cmds = append(cmds, cmd)
		cmds = append(cmds, notifications.CreateNotificationMsg("Fetching assigned issues", notifications.Info, notifications.Bar))
		return m, tea.Batch(cmds...)
	}

	//return m, tea.Batch(cmds...)
}

func (m ChangeAssigneeModel) GetNavTo() navigation.ScreenId {
	return m.NavTo
}

func (m ChangeAssigneeModel) SetNavTo(navTo navigation.ScreenId) navigation.ScreenModel {
	m.NavTo = navTo
	return m
}

func (m ChangeAssigneeModel) GetContext() navigation.Context {

	return m.Context
}

func (m ChangeAssigneeModel) SetContext(context navigation.Context) navigation.ScreenModel {
	m.Context = context
	return m
}
