package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.stacksList.FilterState() == list.Filtering {
		m.stacksList, cmd = m.stacksList.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.message = ""
	m.errorMessage = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			fs := m.stacksList.FilterState()
			if fs == list.Filtering || fs == list.FilterApplied {
				m.stacksList.ResetFilter()
			} else if m.stacksFilter != stacksFilterAll {
				m.stacksList.SetItems(m.stacksListReserve)
				m.stacksList.Title = "stacks"
				m.stacksFilter = stacksFilterAll
				m.stacksList.Styles.Title.Background(lipgloss.Color(stackListColor))
			} else {
				return m, tea.Quit
			}
		case "o":
			if m.stacksFilter != stacksFilterOuttaSync {
				filteredItems := make([]list.Item, 0)
				for _, item := range m.stacksListReserve {
					stackItem, ok := item.(Stack)
					if ok {
						if stackItem.FetchStatus == StatusFetched && stackItem.OuttaSync {
							filteredItems = append(filteredItems, item)
						}
					}
				}
				m.stacksList.SetItems(filteredItems)
				m.stacksList.Title = "stacks (outtasync)"
				m.stacksFilter = stacksFilterOuttaSync
				m.stacksList.Styles.Title.Background(lipgloss.Color(outtaSyncColor))

			}
		case "i":
			if m.stacksFilter != stacksFilterInSync {
				filteredItems := make([]list.Item, 0)
				for _, item := range m.stacksListReserve {
					stackItem, ok := item.(Stack)
					if ok {
						if stackItem.FetchStatus == StatusFetched && !stackItem.OuttaSync {
							filteredItems = append(filteredItems, item)
						}
					}
				}
				m.stacksList.SetItems(filteredItems)
				m.stacksList.Title = "stacks (in sync)"
				m.stacksFilter = stacksFilterInSync
				m.stacksList.Styles.Title.Background(lipgloss.Color(inSyncColor))
			}
		case "e":
			if m.stacksFilter != stacksFilterErr {
				filteredItems := make([]list.Item, 0)
				for _, item := range m.stacksListReserve {
					stackItem, ok := item.(Stack)
					if ok {
						if stackItem.Err != nil {
							filteredItems = append(filteredItems, item)
						}
					}
				}
				m.stacksList.SetItems(filteredItems)
				m.stacksList.Title = "stacks (with errors)"
				m.stacksFilter = stacksFilterErr
				m.stacksList.Styles.Title.Background(lipgloss.Color(errorColor))
			}
		}
	case tea.WindowSizeMsg:
		_, h1 := stackListStyle.GetFrameSize()
		m.stacksList.SetHeight(msg.Height - h1 - 2)
	case CheckStackStatus:
		msg.stack.FetchStatus = StatusFetching
		m.stacksList.SetItem(msg.index, msg.stack)
		m.stacksListReserve = m.stacksList.Items()
		return m, getCFTemplateBody(m.awsConfigs[GetAWSConfigKey(msg.stack)], msg.index, msg.stack)
	case TemplateFetchedMsg:
		if msg.err != nil {
			msg.stack.Err = msg.err
			msg.stack.FetchStatus = StatusFailure
			m.stacksList.SetItem(msg.index, msg.stack)
			m.resultMap[msg.index] = stackResultErr
		} else {
			msg.stack.OuttaSync = true
			msg.stack.FetchStatus = StatusFetched
			msg.stack.OuttaSync = msg.outtaSync
			switch msg.outtaSync {
			case true:
				m.resultMap[msg.index] = stackResultOuttaSync
			case false:
				m.resultMap[msg.index] = stackResultInSync
			}
			msg.stack.Template = msg.template
			msg.stack.Err = nil
			m.stacksList.SetItem(msg.index, msg.stack)
		}

		// recompute outtasync and error numbers
		m.outtaSyncNum = 0
		m.errorNum = 0
		for _, v := range m.resultMap {
			if v == stackResultOuttaSync {
				m.outtaSyncNum++
			} else if v == stackResultErr {
				m.errorNum++
			}
		}
		m.stacksListReserve = m.stacksList.Items()
	case ShowFileFinished:
		if msg.err != nil {
			m.errorMessage = fmt.Sprintf("Error showing file: %s", Trim(msg.err.Error(), 50))
		}
		return m, tea.Batch(cmds...)
	case CredentialsRefreshedMsg:
		if msg.err != nil {
			m.errorMessage = "Error refreshing credentials"
		} else {
			m.message = "Credentials Refreshed"
		}
	case hideHelpMsg:
		m.showHelp = false
	}

	switch m.state {
	case cfStacksList:
		m.stacksList, cmd = m.stacksList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
