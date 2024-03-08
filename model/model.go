package model

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type stateView uint

const (
	cfStacksList stateView = iota
)

type awsConfig struct {
	config aws.Config
	err    error
}

type model struct {
	awsConfigs     map[string]awsConfig
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
