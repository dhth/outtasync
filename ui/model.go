package ui

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type stateView uint

const (
	cfStacksList stateView = iota
)

type AwsConfig struct {
	Config aws.Config
	Err    error
}

type AwsCFClient struct {
	Client *cloudformation.Client
	Err    error
}

type model struct {
	awsConfigs     map[string]AwsConfig
	state          stateView
	stacksList     list.Model
	outtaSyncMap   map[int]int
	checkOnStart   bool
	message        string
	errorMessage   string
	terminalHeight int
	terminalWidth  int
	err            error
	outtaSyncNum   uint
	errorNum       uint
}

func (m model) Init() tea.Cmd {

	var cmds []tea.Cmd
	if m.checkOnStart {
		for i, stack := range m.stacksList.Items() {
			if st, ok := stack.(Stack); ok {
				cmds = append(cmds, StackChosen(i, st))
			}
		}
	}
	return tea.Batch(cmds...)
}
