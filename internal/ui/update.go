package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/outtasync/internal/utils"
)

const (
	numThrottledCmdsUpperLimit = 3
	numSyncCallsUpperLimit     = 30
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
			switch m.activePane {
			case errorDetailsPane:
				m.activePane = stacksList
			case result:
				m.activePane = stacksList
			case stacksList:
				return m, tea.Quit
			default:
				m.goBack()
			}
		case "enter", "s":
			syncCmd, ok := m.getCmdForTemplateCheck()
			if ok {
				cmds = append(cmds, syncCmd)
			}
		case "S":
			cmds = append(cmds, m.getCmdsForTemplateCheck()...)
		case "d":
			driftCmd, ok := m.getCmdForDriftCheck()
			if ok {
				cmds = append(cmds, driftCmd)
			}
		case "D":
			driftCmds := m.getCmdsForDriftCheck()
			if len(m.throttledCmds)+len(driftCmds) > numThrottledCmdsUpperLimit {
				m.throttledCmds = append(m.throttledCmds, driftCmds...)
			} else {
				cmds = append(cmds, driftCmds...)
			}
		case "ctrl+s":
			var si stackItem
			var ok bool

			switch m.activePane {
			case stacksList:
				si, ok = m.stacksList.SelectedItem().(stackItem)
			case codeMismatchStacksList:
				si, ok = m.codeMismatchStacksList.SelectedItem().(stackItem)
			case driftedStacksList:
				si, ok = m.driftedStacksList.SelectedItem().(stackItem)
			case erroredStacksList:
				si, ok = m.erroredStacksList.SelectedItem().(stackItem)
			default:
				break
			}

			if !ok {
				break
			}

			if si.stack.TemplatePath == nil {
				break
			}

			switch si.syncCheckStatus {
			case syncStatusChecked:
				switch si.outtaSync {
				case true:
					cmds = append(cmds, showDiff(si.templateCode, si.actualTemplate))
				case false:
					cmds = append(cmds, showTemplate(si.templateCode))
				}
			}
		case "tab":
			m.goForward()
		case "shift+tab":
			m.goBack()
		case "1":
			m.activePane = stacksList
		case "2":
			m.showOuttaSyncStacks()
		case "3":
			m.showDriftedStacks()
		case "4":
			m.showStacksWithErrors()
		case "ctrl+e":
			switch m.activePane {
			case stacksList:
				var si stackItem
				var ok bool

				switch m.activePane {
				case stacksList:
					si, ok = m.stacksList.SelectedItem().(stackItem)
				case codeMismatchStacksList:
					si, ok = m.codeMismatchStacksList.SelectedItem().(stackItem)
				case driftedStacksList:
					si, ok = m.driftedStacksList.SelectedItem().(stackItem)
				case erroredStacksList:
					si, ok = m.erroredStacksList.SelectedItem().(stackItem)
				default:
					break
				}

				if !ok {
					break
				}

				if si.syncErr == nil && si.driftErr == nil {
					m.message = "no errors for selected stack"
					break
				}

				var errorDetails string
				if si.syncErr != nil {
					errorDetails += fmt.Sprintf(`%s

%s

`, errorDetailsHeadingStyle.Render("Sync Error"), si.syncErr.Error())
				}

				if si.driftErr != nil {
					errorDetails += fmt.Sprintf(`%s

%s
`, errorDetailsHeadingStyle.Render("Drift Check Error"), si.driftErr.Error())
				}

				// TODO: not sure why the viewport is not wrapping content by default
				// this is a workaround for that
				m.stackErrorVP.SetContent(lipgloss.NewStyle().Width(m.terminalWidth - 4).Render(errorDetails))
				m.activePane = errorDetailsPane
			case errorDetailsPane:
				m.activePane = stacksList
			}
		}
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		w, h := stackListStyle.GetFrameSize()
		m.stacksList.SetWidth(msg.Width - w - 2)
		m.stacksList.SetHeight(msg.Height - h - 2)
		m.codeMismatchStacksList.SetWidth(msg.Width - w - 2)
		m.codeMismatchStacksList.SetHeight(msg.Height - h - 2)
		m.driftedStacksList.SetWidth(msg.Width - w - 2)
		m.driftedStacksList.SetHeight(msg.Height - h - 2)
		m.erroredStacksList.SetWidth(msg.Width - w - 2)
		m.erroredStacksList.SetHeight(msg.Height - h - 2)

		if !m.resultVPReady {
			m.resultVP = viewport.New(msg.Width-4, msg.Height-5)
			m.resultVPReady = true
		} else {
			m.resultVP.Width = msg.Width - 4
			m.resultVP.Height = msg.Height - 5
		}

		if !m.stackErrorVPReady {
			m.stackErrorVP = viewport.New(msg.Width-4, msg.Height-5)
			m.stackErrorVPReady = true
		} else {
			m.stackErrorVP.Width = msg.Width - 4
			m.stackErrorVP.Height = msg.Height - 5
		}
	case DriftCheckUpdated:
		si, ok := m.stacksList.Items()[msg.index].(stackItem)
		if !ok {
			break
		}

		si.driftErr = msg.result.Err
		si.driftOutput = msg.result.Output
		si.driftCheckStatus = driftChecked
		m.stacksList.SetItem(msg.index, si)
		m.recomputeStats()
		if msg.throttled && m.throttledCmdsInProgress > 0 {
			m.throttledCmdsInProgress--
		}

	case TemplateFetchedMsg:
		items := m.stacksList.Items()
		if msg.index >= len(items) {
			break
		}

		si, ok := items[msg.index].(stackItem)
		if !ok {
			break
		}

		if msg.err != nil {
			si.syncCheckStatus = syncStatusCheckFailed
		} else {
			si.outtaSync = true
			si.syncCheckStatus = syncStatusChecked
			si.outtaSync = msg.mismatch
			si.templateCode = msg.templateCode
			si.actualTemplate = msg.actualTemplate
		}

		si.syncErr = msg.err
		m.stacksList.SetItem(msg.index, si)

		m.recomputeStats()
		if msg.throttled && m.throttledCmdsInProgress > 0 {
			m.throttledCmdsInProgress--
		}
	case ShowFileFinished:
		if msg.err != nil {
			m.errorMessage = fmt.Sprintf("Error showing file: %s", utils.Trim(msg.err.Error(), 50))
		}
	case hideHelpMsg:
		m.showHelp = false
	}

	switch m.activePane {
	case stacksList:
		m.stacksList, cmd = m.stacksList.Update(msg)
		cmds = append(cmds, cmd)
	case codeMismatchStacksList:
		m.codeMismatchStacksList, cmd = m.codeMismatchStacksList.Update(msg)
		cmds = append(cmds, cmd)
	case driftedStacksList:
		m.driftedStacksList, cmd = m.driftedStacksList.Update(msg)
		cmds = append(cmds, cmd)
	case erroredStacksList:
		m.erroredStacksList, cmd = m.erroredStacksList.Update(msg)
		cmds = append(cmds, cmd)
	case errorDetailsPane:
		m.stackErrorVP, cmd = m.stackErrorVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	if len(m.throttledCmds) > 0 && m.throttledCmdsInProgress < numThrottledCmdsUpperLimit {
		if len(m.throttledCmds) < numThrottledCmdsUpperLimit {
			cmds = append(cmds, m.throttledCmds...)
			m.throttledCmdsInProgress += len(m.throttledCmds)
			m.throttledCmds = make([]tea.Cmd, 0)
		} else {
			numToFetch := numThrottledCmdsUpperLimit - m.throttledCmdsInProgress
			cmds = append(cmds, m.throttledCmds[:numToFetch]...)
			m.throttledCmds = m.throttledCmds[numToFetch:]
			m.throttledCmdsInProgress += numToFetch
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) getCmdForTemplateCheck() (tea.Cmd, bool) {
	if m.activePane != stacksList {
		return nil, false
	}

	si, ok := m.stacksList.SelectedItem().(stackItem)
	if !ok {
		return nil, false
	}

	if si.stack.TemplatePath == nil {
		return nil, false
	}

	index := m.stacksList.Index()
	si.syncCheckStatus = syncStatusInProgress
	si.syncErr = nil
	m.stacksList.SetItem(index, si)
	remoteCallHeaders := append(m.remoteCallHeaders, si.stack.TemplateRemoteCallHeaders...)
	return getCFTemplateBody(
		m.cfClients[si.stack.AWSConfigKey()],
		index,
		si.stack.Name,
		si.stack.Key(),
		*si.stack.TemplatePath,
		remoteCallHeaders,
		false), true
}

func (m *Model) getCmdsForTemplateCheck() []tea.Cmd {
	//nolint:prealloc
	var cmds []tea.Cmd

	if m.activePane != stacksList {
		return cmds
	}

	for i, li := range m.stacksList.Items() {
		si, ok := li.(stackItem)
		if !ok {
			continue
		}

		if si.stack.TemplatePath == nil {
			continue
		}

		si.syncErr = nil
		si.syncCheckStatus = syncStatusInProgress
		m.stacksList.SetItem(i, si)
		remoteCallHeaders := append(m.remoteCallHeaders, si.stack.TemplateRemoteCallHeaders...)
		cmds = append(cmds, getCFTemplateBody(
			m.cfClients[si.stack.AWSConfigKey()],
			i,
			si.stack.Name,
			si.stack.Key(),
			*si.stack.TemplatePath,
			remoteCallHeaders,
			true))
	}

	return cmds
}

func (m *Model) getCmdForDriftCheck() (tea.Cmd, bool) {
	if m.activePane != stacksList {
		return nil, false
	}

	si, ok := m.stacksList.SelectedItem().(stackItem)
	if !ok {
		return nil, false
	}

	cfClient, ok := m.cfClients[si.stack.AWSConfigKey()]
	if !ok {
		return nil, false
	}
	index := m.stacksList.Index()
	si.driftCheckStatus = driftCheckInProgress
	si.driftErr = nil
	m.stacksList.SetItem(index, si)
	return checkStackDrift(cfClient, index, si, false), true
}

func (m *Model) getCmdsForDriftCheck() []tea.Cmd {
	var cmds []tea.Cmd
	if m.activePane != stacksList {
		return cmds
	}

	for i, stack := range m.stacksList.Items() {
		if si, ok := stack.(stackItem); ok {
			awsCfg, ok := m.cfClients[si.stack.AWSConfigKey()]
			if !ok {
				continue
			}

			si.driftErr = nil
			si.driftCheckStatus = driftCheckInProgress
			m.stacksList.SetItem(i, si)
			cmds = append(cmds, checkStackDrift(awsCfg, i, si, true))
		}
	}

	return cmds
}

func (m *Model) showOuttaSyncStacks() {
	if m.activePane == codeMismatchStacksList {
		return
	}

	filteredItems := make([]list.Item, 0)
	for _, li := range m.stacksList.Items() {
		si, ok := li.(stackItem)
		if !ok {
			continue
		}

		if si.syncCheckStatus == syncStatusChecked && si.outtaSync {
			filteredItems = append(filteredItems, si)
		}
	}
	m.codeMismatchStacksList.SetItems(filteredItems)
	m.activePane = codeMismatchStacksList
}

func (m *Model) showDriftedStacks() {
	if m.activePane == driftedStacksList {
		return
	}

	filteredItems := make([]list.Item, 0)
	for _, li := range m.stacksList.Items() {
		si, ok := li.(stackItem)
		if !ok {
			continue
		}

		if si.hasDrifted() {
			filteredItems = append(filteredItems, si)
		}
	}
	m.driftedStacksList.SetItems(filteredItems)
	m.activePane = driftedStacksList
}

func (m *Model) showStacksWithErrors() {
	if m.activePane == erroredStacksList {
		return
	}

	filteredItems := make([]list.Item, 0)
	for _, li := range m.stacksList.Items() {
		si, ok := li.(stackItem)
		if !ok {
			continue
		}
		if si.syncErr != nil || si.driftErr != nil {
			filteredItems = append(filteredItems, si)
		}
	}
	m.erroredStacksList.SetItems(filteredItems)
	m.activePane = erroredStacksList
}

func (m *Model) recomputeStats() {
	m.outtaSyncNum = 0
	m.errorNum = 0
	m.driftedNum = 0

	for _, li := range m.stacksList.Items() {
		si, ok := li.(stackItem)
		if !ok {
			continue
		}

		if si.syncErr != nil {
			m.errorNum++
		} else if si.outtaSync {
			m.outtaSyncNum++
		}

		if si.driftErr != nil {
			m.errorNum++
		} else if si.hasDrifted() {
			m.driftedNum++
		}
	}
}

func (m *Model) goForward() {
	switch m.activePane {
	case stacksList:
		m.showOuttaSyncStacks()
	case codeMismatchStacksList:
		m.showDriftedStacks()
	case driftedStacksList:
		m.showStacksWithErrors()
	case erroredStacksList:
		m.activePane = stacksList
	}
}

func (m *Model) goBack() {
	switch m.activePane {
	case stacksList:
		m.showStacksWithErrors()
	case codeMismatchStacksList:
		m.activePane = stacksList
	case driftedStacksList:
		m.showOuttaSyncStacks()
	case erroredStacksList:
		m.showDriftedStacks()
	}
}
