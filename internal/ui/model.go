package ui

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type stateView uint

const (
	cfStacksList stateView = iota
)

type stackFilter uint

const (
	stacksFilterAll stackFilter = iota
	stacksFilterErr
	stacksFilterOuttaSync
	stacksFilterInSync
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
	awsConfigs        map[string]AwsConfig
	state             stateView
	stacksList        list.Model
	stacksListReserve map[string]Stack
	stacksFilter      stackFilter
	checkOnStart      bool
	message           string
	errorMessage      string
	outtaSyncNum      uint
	errorNum          uint
	showHelp          bool
}

func (m model) Init() tea.Cmd {

	var cmds []tea.Cmd

	cmds = append(cmds, hideHelp(time.Minute*1))

	if m.checkOnStart {
		for i, stack := range m.stacksList.Items() {
			if st, ok := stack.(Stack); ok {
				cmds = append(cmds, StackChosen(i, st))
			}
		}
	}
	return tea.Batch(cmds...)
}
