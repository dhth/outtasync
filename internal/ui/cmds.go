package ui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/outtasync/internal/aws"
)

func StackChosen(index int, si stackItem) tea.Cmd {
	return func() tea.Msg {
		return CheckStackStatus{index, si}
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

func showDiff(si stackItem) tea.Cmd {
	c := exec.Command("bash", "-c",
		fmt.Sprintf("cat << 'EOF' | git diff --dst-prefix='Actual Cloudformation stack' --no-index -- %s -\n%s\nEOF",
			si.stack.Local,
			si.stack.Template,
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

func getCFTemplateBody(awsConfig aws.Config, index int, stackItem stackItem) tea.Cmd {
	return func() tea.Msg {
		stackSyncStatus := aws.CheckStackSyncStatus(awsConfig, stackItem.stack)

		if awsConfig.Err != nil {
			return TemplateFetchedMsg{index, stackItem, "", false, awsConfig.Err}
		}
		return TemplateFetchedMsg{index, stackItem, stackSyncStatus.TemplateBody, stackSyncStatus.Outtasync, stackSyncStatus.Err}
	}
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return hideHelpMsg{}
	})
}
