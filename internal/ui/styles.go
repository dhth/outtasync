package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	stackListColor         = "#fe8019"
	modeColor              = "#b8bb26"
	fetchingColor          = "#ebdbb2"
	inSyncColor            = "#b8bb26"
	outtaSyncColor         = "#fb4934"
	errorColor             = "#928374"
	helpMsgColor           = "#83a598"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	baseListStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	stackListStyle = baseListStyle

	modeStyle = baseStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(modeColor))

	driftStatusStyle = baseStyle.
				Bold(true).
				Align(lipgloss.Center).
				Width(12)

	fetchingStyle = driftStatusStyle.
			Background(lipgloss.Color(fetchingColor))

	insSyncStyle = driftStatusStyle.
			Background(lipgloss.Color(inSyncColor))

	outtaSyncStyle = driftStatusStyle.
			Background(lipgloss.Color(outtaSyncColor))

	errorStyle = driftStatusStyle.
			Background(lipgloss.Color(errorColor))

	msgStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Bold(true)

	outtaSyncMsgStyle = msgStyle.
				Foreground(lipgloss.Color(outtaSyncColor))

	errorMsgStyle = msgStyle.
			Foreground(lipgloss.Color(errorColor))

	helpMsgStyle = baseStyle.
			Bold(true).
			PaddingLeft(2).
			Foreground(lipgloss.Color(helpMsgColor))

	tagColors = []string{
		"#d3869b",
		"#8ec07c",
		"#fabd2f",
		"#83a598",
		"#48cae4",
		"#ff99ac",
		"#ff5c8a",
		"#e0aaff",
	}
	tagStyle = func(tag string) lipgloss.Style {
		h := fnv.New32()
		h.Write([]byte(tag))
		hash := h.Sum32()

		color := tagColors[int(hash)%len(tagColors)]

		st := lipgloss.NewStyle().
			PaddingRight(1).
			Foreground(lipgloss.Color(color))

		return st
	}
)
