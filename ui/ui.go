package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/outtasync/model"
)

func RenderUI(stacks []model.Stack) {
	p := tea.NewProgram(model.InitialModel(stacks), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}
