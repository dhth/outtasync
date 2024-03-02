package model

import (
	"github.com/charmbracelet/bubbles/list"
)

func InitialModel(stacks []Stack) model {
	stackItems := make([]list.Item, 0, len(stacks))
	for _, stack := range stacks {
		stackItems = append(stackItems, stack)
	}
	var appDelegateKeys = newAppDelegateKeyMap()
	appDelegate := newAppItemDelegate(appDelegateKeys)

	m := model{
		stacksList: list.New(stackItems, appDelegate, listPadding, 0),
	}
	m.stacksList.Title = "Stacks"
	m.stacksList.SetFilteringEnabled(false)

	return m
}
