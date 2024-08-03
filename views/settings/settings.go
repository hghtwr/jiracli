package settings

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hghtwr/jiracli/customHelp"
)

func View() string {

	help := customHelp.CreateDefaultHelp()
	return help.View(CreateCustomKeyMap())

}

func Update(msg tea.Msg) tea.Cmd {
	return nil
}




func CreateCustomKeyMap() *customHelp.DefaultKeyMap {
	return customHelp.CreateDefaultKeyMap([][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("z"),
				key.WithHelp("z", "Do something..."),
			),
	},
})
}
