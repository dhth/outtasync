package model

import (
	"github.com/charmbracelet/bubbles/list"
)

func InitialModel(stacks []Stack) model {
	stackItems := make([]list.Item, 0, len(stacks))
	awsCfgs := make(map[string]awsConfig)

	seen := make(map[string]bool)
	for _, stack := range stacks {
		stackItems = append(stackItems, stack)
		configKey := getAWSConfigKey(stack)
		if !seen[configKey] {
			cfg, err := getAWSConfig(stack.AwsProfile, stack.AwsRegion)
			awsCfgs[configKey] = awsConfig{cfg, err}
			seen[configKey] = true
		}
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

	return m
}
