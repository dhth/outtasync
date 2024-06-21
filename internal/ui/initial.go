package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(stacks []Stack, awsCfgs map[string]AwsConfig, checkOnStart bool) model {
	stackItems := make([]list.Item, len(stacks))

	stackReserve := make(map[string]Stack)
	for i, stack := range stacks {
		stackItems[i] = stack
		stackReserve[stack.key()] = stack
	}

	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)

	m := model{
		awsConfigs:   awsCfgs,
		stacksList:   list.New(stackItems, appDelegate, listWidth, 0),
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
