package changeAssignee

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
	form *huh.Form
	formModel tea.Model
	assigneeSearchResult []jira.User
	searchField *huh.Input
	blockNavigation bool
	//searchFieldValue string
	//userSelect *huh.Select[selectValue]
	userSelect *huh.Select[string]
	submitButton *huh.Confirm

}
type DetailFields struct {
	Assignee 	string
}

type selectValue struct {
	displayName string
	accountId string
	email string
}

var (
	searchFieldValue string
)

func CreateInitModel() ChangeAssigneeModel{
	m := ChangeAssigneeModel{
		NavTo: navigation.ChangeAssigneeView,
		assigneeSearchResult: nil,
		blockNavigation: false,
	}
	m.createForm()

	return m
}

func (m ChangeAssigneeModel) Init() tea.Cmd {
	return nil
}

func (m ChangeAssigneeModel) View() string {
	headerStyle := lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.ThickBorder()).Width(layout.GetWidthFraction(24))
	var headerContent  []string
	var assigneeContent []string
	var assigneeForm = m.form.View()

	assigneeContent = append(assigneeContent, layout.Style.DetailsFieldTitleStyle.Render("Current Assignee"))

	if m.issue != nil {
		headerContent = append(headerContent, layout.Style.HeaderStyle.Render(m.issue.Key))
		headerContent = append(headerContent, m.issue.Fields.Summary)

		assigneeContent = append(assigneeContent, layout.Style.DetailsFieldValueStyle.Render(m.issue.Fields.Assignee.DisplayName))
	}

	help := customHelp.CreateDefaultHelp()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left, headerContent...)),
			lipgloss.JoinHorizontal(lipgloss.Left, assigneeContent...),
			"\n\n",
			assigneeForm,
		help.View(customHelp.CreateDefaultKeyMap(nil)))
}

func (m ChangeAssigneeModel) Update(msg tea.Msg) (navigation.ScreenModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
  m.blockNavigation = true
	switch msg := msg.(type) {
	case tea.KeyMsg:
		 if msg.String() == "esc" {
			m.blockNavigation = false
		}
			m.formModel, cmd = m.form.Update(msg)
			cmds = append(cmds, cmd)

			searchField := m.searchField.GetValue() // As long as the form is not completed, we need to use this hack to get the value :/

			if str, ok := searchField.(string); ok {
				if len(str) > 3{
					cmd := jiraApi.SearchAssignees(str)
					cmds = append(cmds, cmd)
				}
			}
			return m, tea.Batch(cmds...)

	case jiraApi.IssueDetailResponse:
		m.issue = msg.Issue
		//Now fetch tasks that have this task as parent (different from subtasks, e.g. for stories in epics)
		return m, tea.Batch(cmds...)

	case cursor.BlinkMsg:
		return m, tea.Batch(cmds...)
	case jiraApi.SearchAssigneesResponse:
		m.assigneeSearchResult = msg.Assignees
		return m, tea.Batch(cmds...)
	default:
		cmd := jiraApi.FetchIssueDetails(m.Context.IssueId)
		cmds = append(cmds, cmd)
		cmds = append(cmds, notifications.CreateNotificationMsg("Fetching assigned issues", notifications.Info, notifications.Bar))
		return m, tea.Batch(cmds...)

	}

	//return m, tea.Batch(cmds...)
}

func (m *ChangeAssigneeModel) createForm() *huh.Form {

	m.searchField = huh.NewInput().Key("search").Title("Search user").Prompt("Username? ").Value(&searchFieldValue)

	m.userSelect = huh.NewSelect[string]().Key("user").TitleFunc(
		func() string {
			return searchFieldValue
		}, &searchFieldValue,
	).OptionsFunc(
		func() []huh.Option[string]{

			var options []huh.Option[string]

			options = append(options, huh.NewOption(
				"test",
				"test",
			),
			)
			return options
		}, &searchFieldValue,



		/*func() []huh.Option[selectValue]{

			var options []huh.Option[selectValue]
				options = append(options, huh.NewOption(
					"test",
						selectValue{
							displayName: "test",
							accountId: "test",
							email: "test",
						},
					),
				)
			return options
		}, &m.searchFieldValue,
	*/
		)

	m.submitButton = huh.NewConfirm().Affirmative("Submit").Validate(
		func(v bool) error {
			if !v {
				return fmt.Errorf("please confirm")
			}
			return nil
		},
	).Negative("Discard").Validate(
		func(v bool) error {
			if !v {
				return fmt.Errorf("Discard")
			}
			return nil
		},
	)

	m.searchField.Focus()


	form := huh.NewForm(huh.NewGroup(m.searchField, m.userSelect, m.submitButton))
	m.form = form

	return form
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

func (m ChangeAssigneeModel) GetBlockNavigation() bool {
	return m.blockNavigation
}

func (m ChangeAssigneeModel) SetBlockNavigation(block bool) navigation.ScreenModel {
	m.blockNavigation = block
	return m
}