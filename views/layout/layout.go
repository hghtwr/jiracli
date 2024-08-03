package layout

import "golang.org/x/term"

func GetWidthFraction(multiplier int) int {
	return GetTerminalWidthWidth() / 24 * multiplier
}

func GetTerminalWidthWidth() int {
	terminalWidth, _, _ := term.GetSize(0)
	return terminalWidth
}