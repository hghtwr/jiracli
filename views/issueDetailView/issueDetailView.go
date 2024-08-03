package issueDetailView

import (
	"fmt"
	"reflect"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hghtwr/jiracli/customHelp"
	"github.com/hghtwr/jiracli/jiraApi"
	"github.com/hghtwr/jiracli/navigation"
	"github.com/hghtwr/jiracli/notifications"
	"github.com/hghtwr/jiracli/views/layout"
)


type IssueDetailViewModel struct {
	NavTo 		navigation.ScreenId
	Context 	navigation.Context
	issueParent *jira.Issue
	issue 		*jira.Issue
	issueChildIssues []jira.Issue
	selectedTab int
}

const (
	CommentTab = iota
	SubtasksTab
	LinksTab
	ChildIssuesTab
)

type DetailFields struct {
	Type 		string
	Assignee 	string
	Status 		string
	Priority 	string
	Description string
	Reporter 	string
	Parent 		string
}

func CreateInitModel() IssueDetailViewModel{
	return IssueDetailViewModel{
		NavTo: navigation.IssueDetailView,
	}
}

func (m IssueDetailViewModel) Init() tea.Cmd {
	return nil
}

func (m IssueDetailViewModel) createViewTabs() []string {
	var commentBoxTabs []string
	tabs :=  []string{
		fmt.Sprintf("Comments(%d)", len(m.issue.Fields.Comments.Comments)),
		fmt.Sprintf("Subtasks(%d)", len(m.issue.Fields.Subtasks)),
		fmt.Sprintf("Links(%d)", len(m.issue.Fields.IssueLinks)),
		fmt.Sprintf("Child issues(%d)", len(m.issueChildIssues)),
	}

	for i, tab := range tabs {
		if i == m.selectedTab {
			commentBoxTabs = append(commentBoxTabs, layout.Style.ActiveCommentBoxStyle.Render(tab))
		} else {
			commentBoxTabs = append(commentBoxTabs, layout.Style.CommentBoxTabStyle.Render(tab))
		}
	}
	return commentBoxTabs
}

func (m IssueDetailViewModel) createCommentBoxContent() []string {
	var commentBoxContent []string
	comments := m.issue.Fields.Comments.Comments

	for _, comment := range comments {
		commentBoxContent = append(commentBoxContent, comment.Created +  " - " + comment.Author.DisplayName + ": " + comment.Body)
	}
	return commentBoxContent
}

func (m IssueDetailViewModel) createSubtaskBoxContent() []string {
	var subTasks []string
	subtasks := m.issue.Fields.Subtasks
	subtaskStyle := lipgloss.NewStyle()
	for _, subtask := range subtasks {
		if subtask.Fields.Status.StatusCategory.Key == "done" {
			subtaskStyle = subtaskStyle.Strikethrough(true)
		}else{
			subtaskStyle = subtaskStyle.Strikethrough(false)
			subtaskStyle = subtaskStyle.Foreground(lipgloss.Color("15"))
		}
		subTasks = append(subTasks, subtaskStyle.Render(subtask.Key + ": " + subtask.Fields.Summary + " (" + subtask.Fields.Status.Name +")"))
	}
	return subTasks
}

func (m IssueDetailViewModel) createChildIssueBoxContent() []string {
	var childIssuesBoxContent []string
	childIssues := m.issueChildIssues
	childIssueStyle := lipgloss.NewStyle()
	for _, childIssue := range childIssues {
		if childIssue.Fields.Status.StatusCategory.Key == "done" {
			childIssueStyle = childIssueStyle.Strikethrough(true)
		}else{
			childIssueStyle = childIssueStyle.Strikethrough(false)
			childIssueStyle = childIssueStyle.Foreground(lipgloss.Color("15"))
		}
		childIssuesBoxContent = append(childIssuesBoxContent, childIssueStyle.Render(childIssue.Key + ": " + childIssue.Fields.Summary + " (" + childIssue.Fields.Status.Name +")"))
	}
return childIssuesBoxContent
}

func (m IssueDetailViewModel) createLinkBoxContent() []string {
	var linkBoxContent []string
	links := m.issue.Fields.IssueLinks
			linkStyle := lipgloss.NewStyle()
			for _, link := range links {
				var message string
				if link.OutwardIssue != nil {

					if link.OutwardIssue.Fields.Status.StatusCategory.Key == "done"   {
						linkStyle = linkStyle.Strikethrough(true)
					}else{
						linkStyle = linkStyle.Strikethrough(false)
					}

					message = linkStyle.Render("<-- " + link.Type.Outward +  " " + link.OutwardIssue.Key + ": " + link.OutwardIssue.Fields.Summary + " (" + link.OutwardIssue.Fields.Status.Name +")")

				}else if link.InwardIssue != nil {
					if link.InwardIssue.Fields.Status.StatusCategory.Key == "done"  {
						linkStyle = linkStyle.Strikethrough(true)
					}else{
						linkStyle = linkStyle.Strikethrough(false)
					}

					message = linkStyle.Render("--> " + link.Type.Inward + " " +  link.InwardIssue.Key + " : " + link.InwardIssue.Fields.Summary + " (" + link.InwardIssue.Fields.Status.Name +")")

				}
				linkBoxContent = append(linkBoxContent, message)
			}
			return linkBoxContent
}

func (m IssueDetailViewModel) View() string {

	var commentBoxTabs []string
	var commentBoxContent []string

	headerStyle := lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.ThickBorder()).Width(layout.GetWidthFraction(24))
	headerContent := make([]string, 0)

	commentBoxStyle := layout.Style.CommentBoxStyle

	detailsBoxContent := []string {
		layout.Style.SectionTitleStyle.Render("Details"),
		"",

	}

	if m.issue != nil {
		headerContent = append(headerContent, layout.Style.HeaderStyle.Render(m.issue.Key))
		headerContent = append(headerContent, m.issue.Fields.Summary)

		commentBoxTabs = m.createViewTabs()

		fields := DetailFields{
			Type: m.issue.Fields.Type.Name,
			Assignee: m.issue.Fields.Assignee.DisplayName,
			Status: m.issue.Fields.Status.Name,
			Priority: m.issue.Fields.Priority.Name,
			Description: m.issue.Fields.Description,
			Reporter: m.issue.Fields.Reporter.DisplayName,
		}
		if m.issueParent != nil {
			fields.Parent = m.issueParent.Key + ": " + m.issueParent.Fields.Summary
		}
		fieldValues := reflect.ValueOf(fields)
		fieldType := reflect.TypeOf(fields)
		for i := 0; i < fieldValues.NumField(); i++ {
			detailsBoxContent = append(detailsBoxContent, lipgloss.JoinHorizontal(lipgloss.Left, layout.Style.DetailsFieldTitleStyle.Render(fieldType.Field(i).Name + ": "), layout.Style.DetailsFieldValueStyle.Render(fieldValues.Field(i).String())))
		}

		switch m.selectedTab {

		case CommentTab:
			commentBoxContent = m.createCommentBoxContent()
		case SubtasksTab:
			commentBoxContent = m.createSubtaskBoxContent()
		case ChildIssuesTab:
			commentBoxContent = m.createChildIssueBoxContent()
		case LinksTab:
			commentBoxContent = m.createLinkBoxContent()
		}
	}

	help := customHelp.CreateDefaultHelp()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		headerStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left, headerContent...)),
			lipgloss.JoinHorizontal(lipgloss.Left, layout.Style.DetailsBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, detailsBoxContent...))),
			commentBoxStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, commentBoxTabs...)),
			lipgloss.JoinHorizontal(lipgloss.Left, commentBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, commentBoxContent...))),
			"\n",
		help.View(CreateIssueDetailsKeyMap()))

}

func (m IssueDetailViewModel) Update(msg tea.Msg) (navigation.ScreenModel, tea.Cmd) {
	var cmds []tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab":
			if m.selectedTab < 3 {
				m.selectedTab++
			}
		case "shift+tab":
			if m.selectedTab > 0 {
				m.selectedTab--
			}
		case "a":
			m.NavTo = navigation.ChangeAssigneeView
			m.Context = navigation.Context{IssueId: m.issue.Key}
			return m, navigation.RefreshView()
		}
	}
	switch msg := msg.(type) {
	case jiraApi.IssueDetailResponse:
		m.issue = msg.Issue
		if m.issue.Fields.Parent != nil {
			cmd := jiraApi.FetchParentIssueDetails(m.issue.Fields.Parent.Key)
			cmds = append(cmds, cmd)
		}else{
			m.issueParent = nil
		}

		//Now fetch tasks that have this task as parent (different from subtasks, e.g. for stories in epics)
		cmd := jiraApi.FetchChildIssues(m.Context.IssueId)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case jiraApi.IssueParentResponse:
		m.issueParent = msg.Issue
		return m, tea.Batch(cmds...)
	case jiraApi.ChildIssuesResponse:
		m.issueChildIssues = msg.Issues
		return m, tea.Batch(cmds...)
	default:
		cmd := jiraApi.FetchIssueDetails(m.Context.IssueId)
		cmds = append(cmds, cmd)
		cmds = append(cmds, notifications.CreateNotificationMsg("Fetching assigned issues", notifications.Info, notifications.Bar))
		return m, tea.Batch(cmds...)
	}


	//return m, tea.Batch(cmds...)
}

func (m IssueDetailViewModel) GetNavTo() navigation.ScreenId {
	return m.NavTo
}

func (m IssueDetailViewModel) SetNavTo(navTo navigation.ScreenId) navigation.ScreenModel {
	m.NavTo = navTo
	return m
}

func (m IssueDetailViewModel) GetContext() navigation.Context {
	return m.Context
}
func (m IssueDetailViewModel) SetContext(context navigation.Context) navigation.ScreenModel {
	m.Context = context
	return m
}

func CreateIssueDetailsKeyMap() *customHelp.DefaultKeyMap {
	return customHelp.CreateDefaultKeyMap([][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("\t"),
				key.WithHelp("tab", "Next tab"),
	),key.NewBinding(
		key.WithKeys("shift+\t"),
		key.WithHelp("shift+tab", "Previos tab")),
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

