package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(stacks []Stack, awsCfgs map[string]AwsConfig) model {
	stackItems := make([]list.Item, 0, len(stacks))

	for _, stack := range stacks {
		stackItems = append(stackItems, stack)
	}

	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)

	m := model{
		awsConfigs: awsCfgs,
		stacksList: list.New(stackItems, appDelegate, listWidth, 0),
	}
	m.stacksList.Title = "Stacks"
	m.stacksList.SetStatusBarItemName("stack", "stacks")
	m.stacksList.DisableQuitKeybindings()
	m.stacksList.Styles.Title.Background(lipgloss.Color(StackListColor))
	m.stacksList.Styles.Title.Foreground(lipgloss.Color(DefaultBackgroundColor))
	m.stacksList.Styles.Title.Bold(true)

	return m
}
