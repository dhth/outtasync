package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""
	m.errorMessage = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			fs := m.stacksList.FilterState()
			if fs == list.Filtering || fs == list.FilterApplied {
				m.stacksList.ResetFilter()
			} else {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		_, h1 := stackListStyle.GetFrameSize()
		m.stacksList.SetHeight(msg.Height - h1 - 2)
	case CheckStackStatus:
		msg.stack.FetchStatus = StatusFetching
		m.stacksList.SetItem(msg.index, msg.stack)
		return m, GetCFTemplateBody(msg.index, msg.stack)
	case TemplateFetchedMsg:
		if msg.err != nil {
			msg.stack.Err = msg.err
			msg.stack.FetchStatus = StatusFailure
			m.stacksList.SetItem(msg.index, msg.stack)
		} else {
			msg.stack.OuttaSync = true
			msg.stack.FetchStatus = StatusFetched
			msg.stack.OuttaSync = msg.outtaSync
			msg.stack.Template = msg.template
			msg.stack.Err = nil
			m.stacksList.SetItem(msg.index, msg.stack)
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
	}

	switch m.state {
	case cfStacksList:
		m.stacksList, cmd = m.stacksList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
