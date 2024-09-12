package navigation

import tea "github.com/charmbracelet/bubbletea"

type ScreenId int


type Navigation struct {
	NavTo ScreenId
}
type Context struct {
	IssueId string
}
type RefreshViewCmd struct {
}
func RefreshView() tea.Cmd {
	return func() tea.Msg {
		return RefreshViewCmd{}
	}
}
type ScreenModel interface {
	Update(tea.Msg) (ScreenModel, tea.Cmd)
	View() string
	GetNavTo() ScreenId
	SetNavTo(ScreenId) ScreenModel
	GetContext() Context
	SetContext(Context) ScreenModel
	GetBlockNavigation() bool
}
/*type ScreenModel struct {
	ScreenUpdater
	NavTo ScreenId
}*/


const (
	AssignedIssueView ScreenId = iota
	IssueDetailView
	ChangeAssigneeView
	CommentView
	ChangePriorityView
	ChangeStatusView
	SettingsView

)
