package model

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func StackChosen(index int, stack Stack) tea.Cmd {
	return func() tea.Msg {
		return CheckStackStatus{index, stack}
	}
}

func refreshCredentials(cmd string) tea.Cmd {
	cmdEls := strings.Split(cmd, " ")
	c := exec.Command(cmdEls[0], cmdEls[1:]...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return CredentialsRefreshedMsg{err}
		}
		return tea.Msg(CredentialsRefreshedMsg{})
	})
}

func showDiff(stack Stack) tea.Cmd {
	c := exec.Command("bash", "-c",
		fmt.Sprintf("cat << EOF | git diff --dst-prefix='Actual Cloudformation stack' --no-index -- %s -\n%s\nEOF",
			stack.Local,
			stack.Template,
		))
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ShowDiffFinished{err}
		}
		return tea.Msg(ShowDiffFinished{})
	})
}

func showFile(filePath string) tea.Cmd {
	c := exec.Command("bash", "-c",
		fmt.Sprintf("cat %s | less",
			filePath,
		))
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ShowFileFinished{err}
		}
		return tea.Msg(ShowFileFinished{})
	})
}
