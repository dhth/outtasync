package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				var allItems []list.Item
				for _, st := range m.stacksListReserve {
					allItems = append(allItems, st)
				}
				m.stacksList.SetItems(allItems)
				m.stacksList.Title = "stacks"
				m.stacksFilter = stacksFilterAll
				m.stacksList.Styles.Title = m.stacksList.Styles.Title.Background(lipgloss.Color(stackListColor))
			} else {
				return m, tea.Quit
			}
		case "o":
			if m.stacksFilter != stacksFilterOuttaSync {
				filteredItems := make([]list.Item, 0)
				for _, st := range m.stacksListReserve {
					if st.fetchStatus == statusFetched && st.outtaSync {
						filteredItems = append(filteredItems, st)
					}
				}
				m.stacksList.SetItems(filteredItems)
				m.stacksList.Title = "stacks (outtasync)"
				m.stacksFilter = stacksFilterOuttaSync
				m.stacksList.Styles.Title = m.stacksList.Styles.Title.Background(lipgloss.Color(outtaSyncColor))

			}
		case "i":
			if m.stacksFilter != stacksFilterInSync {
				filteredItems := make([]list.Item, 0)
				for _, st := range m.stacksListReserve {
					if st.fetchStatus == statusFetched && !st.outtaSync {
						filteredItems = append(filteredItems, st)
					}
				}
				m.stacksList.SetItems(filteredItems)
				m.stacksList.Title = "stacks (in sync)"
				m.stacksFilter = stacksFilterInSync
				m.stacksList.Styles.Title = m.stacksList.Styles.Title.Background(lipgloss.Color(inSyncColor))
			}
		case "e":
			if m.stacksFilter != stacksFilterErr {
				filteredItems := make([]list.Item, 0)
				for _, st := range m.stacksListReserve {
					if st.err != nil {
						filteredItems = append(filteredItems, st)
					}
				}
				m.stacksList.SetItems(filteredItems)
				m.stacksList.Title = "stacks (with errors)"
				m.stacksFilter = stacksFilterErr
				m.stacksList.Styles.Title = m.stacksList.Styles.Title.Background(lipgloss.Color(errorColor))
			}
		}
	case tea.WindowSizeMsg:
		w, h := stackListStyle.GetFrameSize()
		m.stacksList.SetWidth(msg.Width - w - 2)
		m.stacksList.SetHeight(msg.Height - h - 2)
	case CheckStackStatus:
		msg.stackItem.fetchStatus = statusFetching
		m.stacksList.SetItem(msg.index, msg.stackItem)
		m.stacksListReserve[msg.stackItem.stack.Key()] = msg.stackItem
		return m, getCFTemplateBody(m.awsConfigs[msg.stackItem.stack.AWSConfigKey()], msg.index, msg.stackItem)
	case TemplateFetchedMsg:
		if msg.err != nil {
			msg.stackItem.err = msg.err
			msg.stackItem.fetchStatus = statusFailure
			m.stacksList.SetItem(msg.index, msg.stackItem)
		} else {
			msg.stackItem.outtaSync = true
			msg.stackItem.fetchStatus = statusFetched
			msg.stackItem.outtaSync = msg.outtaSync
			msg.stackItem.stack.Template = msg.template
			msg.stackItem.err = nil
			m.stacksList.SetItem(msg.index, msg.stackItem)
		}

		m.stacksListReserve[msg.stackItem.stack.Key()] = msg.stackItem
		// recompute outtasync and error numbers
		m.outtaSyncNum = 0
		m.errorNum = 0
		for _, st := range m.stacksListReserve {
			if st.err != nil {
				m.errorNum++
			} else if st.outtaSync {
				m.outtaSyncNum++
			}
		}
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
