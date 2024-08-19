package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type delegateKeyMap struct {
	choose             key.Binding
	chooseAll          key.Binding
	refreshCredentials key.Binding
	showDiff           key.Binding
	filterOuttaSync    key.Binding
	filterInSync       key.Binding
	filterErrors       key.Binding
	close              key.Binding
}

func newAppDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("ctrl+f", "enter"),
			key.WithHelp("ctrl+f/enter", "check status"),
		),
		chooseAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "check status for all"),
		),
		refreshCredentials: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh aws credentials"),
		),
		showDiff: key.NewBinding(
			key.WithKeys("ctrl+d", "v"),
			key.WithHelp("ctrl+d/v", "show diff"),
		),
		filterOuttaSync: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "filter outtasync stacks"),
		),
		filterInSync: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "filter in-sync stacks"),
		),
		filterErrors: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "filter stacks with errors"),
		),
		close: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "return to previous page/quit"),
		),
	}
}

type CheckStackStatus struct {
	index     int
	stackItem stackItem
}

func newAppItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(stackListColor)).
		BorderLeftForeground(lipgloss.Color(stackListColor))

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		fs := m.FilterState()
		switch fs {
		case list.FilterApplied:
			switch msgType := msg.(type) {
			case tea.KeyMsg:
				if !key.Matches(msgType, keys.showDiff) {
					return nil
				}
			}
		case list.Filtering:
			return nil
		}
		var si stackItem

		var cmds []tea.Cmd
		index := m.Index()
		if i, ok := m.SelectedItem().(stackItem); ok {
			si = i
		} else {
			return nil
		}

		switch msgType := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msgType, keys.choose):
				return StackChosen(index, si)
			case key.Matches(msgType, keys.chooseAll):
				for i, stack := range m.Items() {
					if st, ok := stack.(stackItem); ok {
						cmds = append(cmds, StackChosen(i, st))
					}
				}
				return tea.Batch(cmds...)
			case key.Matches(msgType, keys.refreshCredentials):
				return refreshCredentials(si.stack.RefreshCommand)
			case key.Matches(msgType, keys.showDiff):
				switch si.fetchStatus {
				case statusFetched:
					switch si.outtaSync {
					case true:
						return showDiff(si)
					case false:
						return showFile(si.stack.Local)
					}
				}
			}
		}
		return nil
	}
	help := []key.Binding{
		keys.choose,
		keys.chooseAll,
		keys.refreshCredentials,
		keys.showDiff,
		keys.filterOuttaSync,
		keys.filterInSync,
		keys.filterErrors,
		keys.close,
	}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	return d
}
