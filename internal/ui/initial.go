package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

func InitialModel(stacks []types.Stack, awsCfgs map[string]aws.Config, checkOnStart bool) Model {
	stackItems := make([]list.Item, len(stacks))

	stackReserve := make(map[string]stackItem)
	for i, stack := range stacks {
		si := stackItem{
			stack: stack,
		}
		stackItems[i] = si
		stackReserve[stack.Key()] = si
	}

	appDelegateKeys := newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)

	m := Model{
		awsConfigs:   awsCfgs,
		stacksList:   list.New(stackItems, appDelegate, defaultListWidth, 0),
		checkOnStart: checkOnStart,
		showHelp:     true,
	}
	m.stacksListReserve = stackReserve
	m.stacksList.Title = "Stacks"
	m.stacksList.SetStatusBarItemName("stack", "stacks")
	m.stacksList.DisableQuitKeybindings()
	m.stacksList.Styles.Title = m.stacksList.Styles.Title.
		Background(lipgloss.Color(stackListColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)

	return m
}
