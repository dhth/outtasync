package model

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	listPadding = 100
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
	footerStr := fmt.Sprintf("%s %s",
		modeStyle.Render("outtasync"),
		errorMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
