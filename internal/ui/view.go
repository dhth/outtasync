package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	listWidth = 140
)

func (m model) View() string {
	var content string
	var footer string

	var statusBar string
	if m.message != "" {
		statusBar = Trim(m.message, 120)
	}
	var errorMsg string
	if m.errorMessage != "" {
		errorMsg = "error: " + Trim(m.errorMessage, 120)
	}

	switch m.state {
	case cfStacksList:
		content = stackListStyle.Render(m.stacksList.View())
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282828")).
		Background(lipgloss.Color("#7c6f64"))

	var outtaSyncMsg string
	var errorCountMsg string

	if m.outtaSyncNum > 0 {
		outtaSyncMsg = outtaSyncMsgStyle.Render(fmt.Sprintf("%d❗", m.outtaSyncNum))
	}
	if m.errorNum > 0 {
		errorCountMsg = errorMsgStyle.Render(fmt.Sprintf("%d 😵", m.errorNum))
	}

	var helpMsg string
	if m.showHelp {
		helpMsg = helpMsgStyle.Render("press ? for help")
	}

	footerStr := fmt.Sprintf("%s%s%s%s  %s",
		modeStyle.Render("outtasync"),
		helpMsg,
		outtaSyncMsg,
		errorCountMsg,
		errorMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
