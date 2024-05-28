package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultBackgroundColor = "#282828"
	StackListColor         = "#fe8019"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#282828"))

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingLeft(1).PaddingBottom(1).Width(listWidth + 10)

	stackListStyle = baseListStyle.Copy()

	modeStyle = baseStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26"))

	driftStatusStyle = baseStyle.Copy().
				Bold(true).
				Align(lipgloss.Center).
				Width(12)

	fetchingStyle = driftStatusStyle.Copy().
			Background(lipgloss.Color("#ebdbb2"))

	insSyncStyle = driftStatusStyle.Copy().
			Background(lipgloss.Color("#b8bb26"))

	outtaSyncStyle = driftStatusStyle.Copy().
			Background(lipgloss.Color("#fb4934"))

	errorStyle = driftStatusStyle.Copy().
			Background(lipgloss.Color("#928374"))

	msgStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Bold(true)

	outtaSyncMsgStyle = msgStyle.Copy().
				Foreground(lipgloss.Color("#fb4934"))

	errorMsgStyle = msgStyle.Copy().
			Foreground(lipgloss.Color("#928374"))
)
