package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(stacks []Stack, awsCfgs map[string]AwsConfig, checkOnStart bool) model {
	stackItems := make([]list.Item, len(stacks))

	resultMap := make(map[int]stackResult)
	for i, stack := range stacks {
		stackItems[i] = stack
		resultMap[i] = stackResultUnchecked
	}

	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)

	m := model{
		awsConfigs:   awsCfgs,
		stacksList:   list.New(stackItems, appDelegate, listWidth, 0),
		checkOnStart: checkOnStart,
		resultMap:    resultMap,
		showHelp:     true,
	}
	m.stacksListReserve = m.stacksList.Items()
	m.stacksList.Title = "Stacks"
	m.stacksList.SetStatusBarItemName("stack", "stacks")
	m.stacksList.DisableQuitKeybindings()
	m.stacksList.Styles.Title.Background(lipgloss.Color(stackListColor))
	m.stacksList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor))
	m.stacksList.Styles.Title.Bold(true)

	return m
}
