package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor  = "#282828"
	stackListColor          = "#fe8019"
	modeColor               = "#b8bb26"
	fetchingColor           = "#ebdbb2"
	inSyncColor             = "#b8bb26"
	unknownDriftStatusColor = "#bdae93"
	driftReasonColor        = "#fabd2f"
	outtaSyncColor          = "#fb4934"
	errorColor              = "#928374"
	errorHeadingColor       = "#fabd2f"
	helpMsgColor            = "#83a598"
)

var (
	baseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(defaultBackgroundColor))

	baseListStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1)

	stackListStyle = baseListStyle

	modeStyle = baseStyle.
			PaddingLeft(1).
			PaddingRight(1).
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(modeColor))

	statusStyle = baseStyle.
			Bold(true).
			Align(lipgloss.Center).
			Width(11)

	fetchingStyle = statusStyle.
			Background(lipgloss.Color(fetchingColor))

	insSyncStyle = statusStyle.
			Background(lipgloss.Color(inSyncColor))

	outtaSyncStyle = statusStyle.
			Background(lipgloss.Color(outtaSyncColor))

	errorStyle = statusStyle.
			Background(lipgloss.Color(errorColor))

	driftCheckInProgressStyle = statusStyle.
					Background(lipgloss.Color(fetchingColor))

	driftedStyle = statusStyle.
			Background(lipgloss.Color(outtaSyncColor))

	notDriftedStyle = statusStyle.
			Background(lipgloss.Color(inSyncColor))

	unknownDriftStatusStyle = statusStyle.
				Background(lipgloss.Color(unknownDriftStatusColor))

	driftReasonStyle = statusStyle.
				Foreground(lipgloss.Color(driftReasonColor))

	driftErrorStyle = statusStyle.
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

	errorViewStyle = baseStyle.
			PaddingLeft(1).
			PaddingRight(1).
			Bold(true).
			Background(lipgloss.Color(errorColor)).
			Align(lipgloss.Left)

	errorDetailsHeadingStyle = baseStyle.
					Foreground(lipgloss.Color(errorHeadingColor))
)
