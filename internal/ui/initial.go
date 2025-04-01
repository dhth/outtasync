package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

func InitialModel(stacks []types.Stack, cfClients map[string]aws.CFClient) Model {
	stackItems := make([]list.Item, len(stacks))

	for i, stack := range stacks {
		si := stackItem{
			stack: stack,
		}
		stackItems[i] = si
	}

	networkCallCmds := make([]tea.Cmd, 0)

	m := Model{
		cfClients:              cfClients,
		stacksList:             list.New(stackItems, newAppItemDelegate(lipgloss.Color(stackListColor)), 0, 0),
		codeMismatchStacksList: list.New(make([]list.Item, 0), newAppItemDelegate(lipgloss.Color(outtaSyncColor)), 0, 0),
		driftedStacksList:      list.New(make([]list.Item, 0), newAppItemDelegate(lipgloss.Color(outtaSyncColor)), 0, 0),
		erroredStacksList:      list.New(make([]list.Item, 0), newAppItemDelegate(lipgloss.Color(errorColor)), 0, 0),
		showHelp:               true,
		throttledCmds:          networkCallCmds,
	}

	m.stacksList.Title = "Stacks"
	m.stacksList.SetStatusBarItemName("stack", "stacks")
	m.stacksList.DisableQuitKeybindings()
	m.stacksList.Styles.Title = m.stacksList.Styles.Title.
		Background(lipgloss.Color(stackListColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)
	m.stacksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.stacksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	m.stacksList.SetShowHelp(false)

	m.codeMismatchStacksList.Title = "outtasync"
	m.codeMismatchStacksList.SetStatusBarItemName("stack", "stacks")
	m.codeMismatchStacksList.DisableQuitKeybindings()
	m.codeMismatchStacksList.Styles.Title = m.codeMismatchStacksList.Styles.Title.
		Background(lipgloss.Color(outtaSyncColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)
	m.codeMismatchStacksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.codeMismatchStacksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	m.codeMismatchStacksList.SetShowHelp(false)

	m.driftedStacksList.Title = "drifted"
	m.driftedStacksList.SetStatusBarItemName("stack", "stacks")
	m.driftedStacksList.DisableQuitKeybindings()
	m.driftedStacksList.Styles.Title = m.driftedStacksList.Styles.Title.
		Background(lipgloss.Color(outtaSyncColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)
	m.driftedStacksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.driftedStacksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	m.driftedStacksList.SetShowHelp(false)

	m.erroredStacksList.Title = "errored"
	m.erroredStacksList.SetStatusBarItemName("stack", "stacks")
	m.erroredStacksList.DisableQuitKeybindings()
	m.erroredStacksList.Styles.Title = m.erroredStacksList.Styles.Title.
		Background(lipgloss.Color(errorColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)
	m.erroredStacksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.erroredStacksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	m.erroredStacksList.SetShowHelp(false)

	return m
}
