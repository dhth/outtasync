package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	}
}

type CheckStackStatus struct {
	index int
	stack Stack
}

func newAppItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(StackListColor)).
		BorderLeftForeground(lipgloss.Color(StackListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle.
		Copy()

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
		var stack Stack

		var cmds []tea.Cmd
		index := m.Index()
		if i, ok := m.SelectedItem().(Stack); ok {
			stack = i
		} else {
			return nil
		}

		switch msgType := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msgType, keys.choose):
				return StackChosen(index, stack)
			case key.Matches(msgType, keys.chooseAll):
				for i, stack := range m.Items() {
					if st, ok := stack.(Stack); ok {
						cmds = append(cmds, StackChosen(i, st))
					}
				}
				return tea.Batch(cmds...)
			case key.Matches(msgType, keys.refreshCredentials):
				return refreshCredentials(stack.RefreshCommand)
			case key.Matches(msgType, keys.showDiff):
				switch stack.FetchStatus {
				case StatusFetched:
					switch stack.OuttaSync {
					case true:
						return showDiff(stack)
					case false:
						return showFile(stack.Local)
					}
				}
			}

		}
		return nil
	}
	help := []key.Binding{keys.choose, keys.chooseAll, keys.refreshCredentials, keys.showDiff}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	return d
}
