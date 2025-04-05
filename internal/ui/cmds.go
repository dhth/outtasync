package ui

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

func checkStackDrift(
	awsConfig aws.CFClient, index int, stackItem stackItem, throttled bool,
) tea.Cmd {
	return func() tea.Msg {
		driftCheckResult := aws.CheckStackDriftStatus(awsConfig, stackItem.stack)
		return DriftCheckUpdated{
			index:     index,
			result:    driftCheckResult,
			throttled: throttled,
		}
	}
}

func showDiff(template, actual string) tea.Cmd {
	tmpFileTemplate, createTmpFileForTemplateErr := os.CreateTemp("", "outtasync-*.yml")
	defer func() {
		_ = tmpFileTemplate.Close()
	}()
	if createTmpFileForTemplateErr != nil {
		return func() tea.Msg {
			return ShowDiffFinished{createTmpFileForTemplateErr}
		}
	}

	_, writeTemplateErr := tmpFileTemplate.WriteString(template)
	if writeTemplateErr != nil {
		return func() tea.Msg {
			return ShowDiffFinished{writeTemplateErr}
		}
	}

	tmpFileActual, createTmpFileForActualErr := os.CreateTemp("", "outtasync-*.yml")
	defer func() {
		_ = tmpFileActual.Close()
	}()
	if createTmpFileForActualErr != nil {
		return func() tea.Msg {
			return ShowDiffFinished{createTmpFileForActualErr}
		}
	}

	_, writeActualErr := tmpFileActual.WriteString(actual)
	if writeActualErr != nil {
		return func() tea.Msg {
			return ShowDiffFinished{writeActualErr}
		}
	}

	c := exec.Command("bash", "-c",
		fmt.Sprintf("git diff --src-prefix='Template ' --dst-prefix='Actual ' --no-index -- %s %s",
			tmpFileTemplate.Name(),
			tmpFileActual.Name(),
		))
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ShowDiffFinished{err}
		}
		return tea.Msg(ShowDiffFinished{})
	})
}

func showTemplate(templateCode string) tea.Cmd {
	tmpFile, err := os.CreateTemp("", "outtasync-*.yml")
	defer func() {
		_ = tmpFile.Close()
	}()
	if err != nil {
		return func() tea.Msg {
			return ShowFileFinished{err}
		}
	}

	_, writeErr := tmpFile.WriteString(templateCode)
	if writeErr != nil {
		return func() tea.Msg {
			return ShowFileFinished{writeErr}
		}
	}
	c := exec.Command("bash", "-c",
		fmt.Sprintf("cat %s | less",
			tmpFile.Name(),
		))
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ShowFileFinished{err}
		}
		return tea.Msg(ShowFileFinished{})
	})
}

func getCFTemplateBody(
	cfClient aws.CFClient,
	index int,
	stackName,
	stackKey,
	templatePath string,
	remoteCallHeaders []types.TemplateRemoteCallHeaders,
	throttled bool,
) tea.Cmd {
	return func() tea.Msg {
		result := aws.CompareStackTemplateCode(cfClient, stackName, stackKey, templatePath, remoteCallHeaders, false)

		if cfClient.Err != nil {
			return TemplateFetchedMsg{
				index:     index,
				throttled: throttled,
				err:       cfClient.Err,
			}
		}
		return TemplateFetchedMsg{
			index:          index,
			templateCode:   result.TemplateCode,
			actualTemplate: result.ActualTemplate,
			mismatch:       result.Mismatch,
			throttled:      throttled,
			err:            result.Err,
		}
	}
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return hideHelpMsg{}
	})
}
