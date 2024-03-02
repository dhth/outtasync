package model

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type stateView uint

const (
	cfStacksList stateView = iota
)

type model struct {
	state          stateView
	stacksList     list.Model
	message        string
	errorMessage   string
	terminalHeight int
	terminalWidth  int
	err            error
}

func (m model) Init() tea.Cmd {
	return nil
}
