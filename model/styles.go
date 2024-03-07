package model

import (
	"github.com/charmbracelet/lipgloss"
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
)
