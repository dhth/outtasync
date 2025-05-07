package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

type Model struct {
	cfClients               map[string]aws.CFClient
	activePane              pane
	lastPane                pane
	stacksList              list.Model
	codeMismatchStacksList  list.Model
	driftedStacksList       list.Model
	erroredStacksList       list.Model
	resultVP                viewport.Model
	resultVPReady           bool
	stackErrorVP            viewport.Model
	stackErrorVPReady       bool
	message                 string
	errorMessage            string
	outtaSyncNum            uint
	errorNum                uint
	driftedNum              uint
	showHelp                bool
	throttledCmds           []tea.Cmd
	throttledCmdsInProgress int
	remoteCallHeaders       []types.RemoteCallHeaders
	helpVP                  viewport.Model
	helpVPReady             bool
	terminalWidth           int
	terminalHeight          int
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, hideHelp(time.Minute*1))

	return tea.Batch(cmds...)
}
