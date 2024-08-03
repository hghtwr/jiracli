package layout

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)


type Styles struct {
	CommentBoxTabStyle       lipgloss.Style
	ActiveCommentBoxStyle    lipgloss.Style
	SectionTitleStyle        lipgloss.Style
	HeaderStyle              lipgloss.Style
	CommentBoxStyle          lipgloss.Style
	DetailsBoxStyle          lipgloss.Style
	DetailsFieldTitleStyle   lipgloss.Style
	DetailsFieldValueStyle   lipgloss.Style
	SubtaskStyle             lipgloss.Style
	ChildIssueStyle          lipgloss.Style
	LinkStyle                lipgloss.Style
}

var Style = &Styles{
	CommentBoxTabStyle:    lipgloss.NewStyle().BorderTop(true).BorderStyle(lipgloss.ThickBorder()).Background(lipgloss.Color("15")).Padding(0, 2).Foreground(lipgloss.Color("0")),
	ActiveCommentBoxStyle: lipgloss.NewStyle().BorderTop(true).BorderStyle(lipgloss.ThickBorder()).Background(lipgloss.Color("14")).Padding(0, 2).Foreground(lipgloss.Color("0")),
	SectionTitleStyle:     lipgloss.NewStyle().Background(lipgloss.Color("14")).Padding(0, 5, 0, 1).Foreground(lipgloss.Color("0")),
	HeaderStyle:           lipgloss.NewStyle().Width(GetWidthFraction(2)).Background(lipgloss.Color("55")).Padding(0, 5, 0, 1).Margin(0, 5, 0, 0),
	CommentBoxStyle:       lipgloss.NewStyle().Width(GetWidthFraction(24)),
	DetailsBoxStyle:       lipgloss.NewStyle().BorderTop(true).Width(GetWidthFraction(24)),
	DetailsFieldTitleStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Padding(0, 0, 0, 1).Width(GetWidthFraction(3)),
	DetailsFieldValueStyle: lipgloss.NewStyle().Padding(0, 0, 0, 1).Width(GetWidthFraction(21)),
	SubtaskStyle:          lipgloss.NewStyle(),
	ChildIssueStyle:       lipgloss.NewStyle(),
	LinkStyle:             lipgloss.NewStyle(),
}
func GetWidthFraction(multiplier int) int {
	return GetTerminalWidthWidth() / 24 * multiplier
}

func GetTerminalWidthWidth() int {
	terminalWidth, _, _ := term.GetSize(0)
	return terminalWidth
}
