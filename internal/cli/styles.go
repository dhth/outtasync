package cli

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	headingColor           = "#83a598"
	defaultBackgroundColor = "#282828"
	inSyncColor            = "#b8bb26"
	outtaSyncColor         = "#fb4934"
	errorColor             = "#928374"
)

var (
	stackNameStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(60)

	stackHeadingStyle = stackNameStyle.
				Bold(true)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Left).
			Width(20)

	inSyncStyle = statusStyle.
			Foreground(lipgloss.Color(inSyncColor))

	stackPositiveResultStyle = stackNameStyle.
					Foreground(lipgloss.Color(inSyncColor))

	stackErrorResultStyle = stackNameStyle.
				Foreground(lipgloss.Color(errorColor))

	stackNegativeResultStyle = stackNameStyle.
					Foreground(lipgloss.Color(outtaSyncColor))

	outtaSyncStyle = statusStyle.
			Foreground(lipgloss.Color(outtaSyncColor))

	errorStyle = statusStyle.
			Foreground(lipgloss.Color(errorColor))
)
