package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/outtasync/internal/utils"
)

func (m Model) View() string {
	var footer string

	var statusBar string
	if m.message != "" {
		statusBar = utils.Trim(m.message, 120)
	}
	var errorMsg string
	if m.errorMessage != "" {
		errorMsg = "error: " + utils.Trim(m.errorMessage, 120)
	}

	var content string
	switch m.activePane {
	case stacksList:
		content = stackListStyle.Render(m.stacksList.View())
	case codeMismatchStacksList:
		content = stackListStyle.Render(m.codeMismatchStacksList.View())
	case driftedStacksList:
		content = stackListStyle.Render(m.driftedStacksList.View())
	case erroredStacksList:
		content = stackListStyle.Render(m.erroredStacksList.View())
	case errorDetailsPane:
		errorViewTitle := errorDetailsTitleStyle.Render("Error(s)")
		if !m.stackErrorVPReady {
			content = vpStyle.Render(lipgloss.JoinVertical(lipgloss.Left, "", errorViewTitle, "", "not ready"))
		} else {
			content = vpStyle.Render(lipgloss.JoinVertical(lipgloss.Left, "", errorViewTitle, "", m.stackErrorVP.View()))
		}
	case helpPane:
		helpTitle := helpViewTitle.Render("Help")
		if !m.helpVPReady {
			content = vpStyle.Render(lipgloss.JoinVertical(lipgloss.Left, "", helpTitle, "", "not ready"))
		} else {
			content = vpStyle.Render(lipgloss.JoinVertical(lipgloss.Left, "", helpTitle, "", m.helpVP.View()))
		}
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282828")).
		Background(lipgloss.Color("#7c6f64"))

	var outtaSyncMsg string
	var syncOrDriftErrMsg string
	var driftedMsg string

	if m.outtaSyncNum > 0 {
		outtaSyncMsg = outtaSyncMsgStyle.Render(fmt.Sprintf("%do", m.outtaSyncNum))
	}
	if m.driftedNum > 0 {
		driftedMsg = outtaSyncMsgStyle.Render(fmt.Sprintf("%dd", m.driftedNum))
	}
	if m.errorNum > 0 {
		syncOrDriftErrMsg = errorMsgStyle.Render(fmt.Sprintf("%de", m.errorNum))
	}

	var helpMsg string
	if m.showHelp {
		helpMsg = helpMsgStyle.Render("press ? for help")
	}

	footerStr := fmt.Sprintf("%s%s%s%s%s  %s",
		modeStyle.Render("outtasync"),
		helpMsg,
		outtaSyncMsg,
		driftedMsg,
		syncOrDriftErrMsg,
		errorMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
