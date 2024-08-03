package notifications

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type NotificationType int
type NotificationMode string

const (
	Info NotificationType = iota
	Warning
	Error
	Success
)

const (
	Tray NotificationMode = "tray"
	Bar NotificationMode = "bar"
)

var notificationIcons =  map[NotificationType]string{
	Success: "✅",
	Error: "❌",
	Warning: "⚠️",
	Info: "ℹ️",
}

type NotificationCmd struct {
	Message string
	Type NotificationType
	Mode NotificationMode
}


func Build(notification *NotificationCmd) string{
	if(notification.Message == ""){
		return ""
	}
	return fmt.Sprintf("%s %s", notificationIcons[notification.Type], notification.Message)
}

func CreateNotificationMsg(message string, notificationType NotificationType, notificationMode NotificationMode) tea.Cmd{
	return func() tea.Msg{
		return NotificationCmd{Message: message, Type: notificationType, Mode: notificationMode}
	}
}

//TO-DO: Check if obsolete.
func ShowNotification(notification string, notificationType NotificationType) tea.Cmd {
	return func() tea.Msg {
		notificationMessage := fmt.Sprintf("%s %s", notificationIcons[notificationType], notification)
		return NotificationCmd{Message: notificationMessage, Type: notificationType}
	}
}