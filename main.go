package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hghtwr/jiracli/navigation"
	"github.com/hghtwr/jiracli/notifications"
	"github.com/hghtwr/jiracli/views/assignedIssues"
	"github.com/hghtwr/jiracli/views/changeAssignee"
	"github.com/hghtwr/jiracli/views/issueDetailView"
)

type MainModel struct {
	activeView          navigation.ScreenId
	notification        notifications.NotificationCmd
	notificationHistory []notifications.NotificationCmd
	history             []navigation.ScreenId
	models              []navigation.ScreenModel
	//assignedIssuesModel assignedIssues.AssignedIssuesModel
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func initialModel() MainModel {

	return MainModel{
		activeView: navigation.AssignedIssueView,
		history: make([]navigation.ScreenId, 0),
		notification: notifications.NotificationCmd{
			Message: "Hello, World!",
			Type:    notifications.Error,
			Mode:    notifications.Tray,
		},
		models: []navigation.ScreenModel{ // This makes update/view much more easier as we can refer to the index of the array from navigation iota. But order needs to be the same as navigation.ScreenId constant!
			assignedIssues.CreateInitModel(),
			issueDetailView.CreateInitModel(),
			changeAssignee.CreateInitModel(),
		},
	}

}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd
	blockNavigation := m.models[m.activeView].GetBlockNavigation()
	switch msg := msg.(type) {
	case notifications.NotificationCmd:
		m.handleNotifications(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {

		case "0":
			if m.activeView != navigation.SettingsView && !blockNavigation {
				m.history = append(m.history, m.activeView)
			}
			m.activeView = navigation.SettingsView
		case "1":
			if m.activeView != navigation.AssignedIssueView && !blockNavigation{
				m.history = append(m.history, m.activeView)
			}
			m.activeView = navigation.AssignedIssueView
		case "esc":
			if len(m.history) > 0 {
				m.models[m.activeView] = m.models[m.activeView].SetNavTo(m.activeView) // Reset the navigation in the current view
				m.activeView = m.history[len(m.history)-1]
				m.history = m.history[:len(m.history)-1]
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	model, cmd := m.models[m.activeView].Update(msg)
	cmdBatch = append(cmdBatch, cmd)
	m.models[m.activeView] = model

	// In case navigation is triggered by subview. We need to reevaluate the blocked State in case it changed
	if m.activeView != m.models[m.activeView].GetNavTo() && !m.models[m.activeView].GetBlockNavigation(){
		m.history = append(m.history, m.activeView)
		navTo := m.models[m.activeView].GetNavTo()
		context := m.models[m.activeView].GetContext()
		m.models[m.activeView] = m.models[m.activeView].SetNavTo(m.activeView) // to not trigger navigation next time we visit this page
		m.activeView = navTo
		m.models[m.activeView] = m.models[m.activeView].SetContext(context)
	}
	return m, tea.Batch(cmdBatch...)

}

func (m MainModel) View() string {

	var mainViewElements []string
	mainViewElements = append(mainViewElements, m.models[m.activeView].View())
	mainViewElements = append(mainViewElements, renderNotifications(&m.notificationHistory, &m.notification))
	return lipgloss.JoinVertical(lipgloss.Left, mainViewElements...)
}

func main() {

	model := initialModel()

	//p := tea.NewProgram(model)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *MainModel) handleNotifications(msg notifications.NotificationCmd) {

	if m.notification.Mode == notifications.Tray {
		m.notificationHistory = append(m.notificationHistory, m.notification)
	}
	if len(m.notificationHistory) > 3 {
		m.notificationHistory = m.notificationHistory[1:]
	}
	m.notification = msg
}

func renderNotifications(notificationHistory *[]notifications.NotificationCmd, notification *notifications.NotificationCmd) string {
	var notificationTray []string

	historyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5e5e5e"))
	barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3e3e3e"))
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f0f0"))

	if notification.Mode == notifications.Bar {
		notificationTray = append(notificationTray, barStyle.Render(notifications.Build(notification)))
		notificationTray = append(notificationTray, "---- Event History----")

	} else {
		notificationTray = append(notificationTray, "")
		notificationTray = append(notificationTray, "---- Event History----")
		notificationTray = append(notificationTray, style.Render(notifications.Build(notification)))

	}

	for i := len(*notificationHistory) - 1; i >= 0; i-- {
		notificationTray = append(notificationTray, historyStyle.Render(notifications.Build(&(*notificationHistory)[i])))
	}

	if len(*notificationHistory) > 3 {
		notificationTray = notificationTray[:3]
	}
	return lipgloss.JoinVertical(lipgloss.Left, notificationTray...)
}
