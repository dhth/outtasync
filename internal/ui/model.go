package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/outtasync/internal/aws"
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

type Model struct {
	awsConfigs        map[string]aws.Config
	state             stateView
	stacksList        list.Model
	stacksListReserve map[string]stackItem
	stacksFilter      stackFilter
	checkOnStart      bool
	message           string
	errorMessage      string
	outtaSyncNum      uint
	errorNum          uint
	showHelp          bool
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, hideHelp(time.Minute*1))

	if m.checkOnStart {
		for i, stack := range m.stacksList.Items() {
			if st, ok := stack.(stackItem); ok {
				cmds = append(cmds, StackChosen(i, st))
			}
		}
	}
	return tea.Batch(cmds...)
}
