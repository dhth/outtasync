package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(stacks []Stack, awsCfgs map[string]AwsConfig, checkOnStart bool) model {
	stackItems := make([]list.Item, 0, len(stacks))

	// -2: error
	// -1: not checked
	//  0: in sync
	//  1: outtasync
	outtaSyncMap := make(map[int]int)
	for i, stack := range stacks {
		stackItems = append(stackItems, stack)
		outtaSyncMap[i] = -1
	}

	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)

	m := model{
		awsConfigs:   awsCfgs,
		stacksList:   list.New(stackItems, appDelegate, listWidth, 0),
		checkOnStart: checkOnStart,
		outtaSyncMap: outtaSyncMap,
	}
	m.stacksList.Title = "Stacks"
	m.stacksList.SetStatusBarItemName("stack", "stacks")
	m.stacksList.DisableQuitKeybindings()
	m.stacksList.Styles.Title.Background(lipgloss.Color(stackListColor))
	m.stacksList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor))
	m.stacksList.Styles.Title.Bold(true)

	return m
}
