package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

type fetchStatus uint

const (
	StatusUnfetched fetchStatus = iota
	StatusFetching
	StatusFetched
	StatusFailure
)

type stackResult uint

const (
	stackResultErr stackResult = iota
	stackResultUnchecked
	stackResultInSync
	stackResultOuttaSync
)

type Stack struct {
	Name           string
	AwsProfile     string
	AwsRegion      string
	Template       string
	Local          string
	Tag            *string
	RefreshCommand string
	FetchStatus    fetchStatus
	OuttaSync      bool
	Err            error
}

func (stack Stack) Title() string {
	return fmt.Sprintf("%s", RightPadTrim(stack.Name, int(float64(listWidth)*0.8)))
}
func (stack Stack) Description() string {
	var status string
	switch stack.FetchStatus {
	case StatusFetched:
		switch stack.OuttaSync {
		case true:
			status = outtaSyncStyle.Render("outta sync")
		case false:
			status = insSyncStyle.Render("in sync")
		}
	case StatusFetching:
		status = fetchingStyle.Render("...")
	case StatusFailure:
		status = errorStyle.Render("error")
	}
	var desc = stack.AwsProfile

	if stack.Err != nil {
		desc = stack.Err.Error()
	}
	return fmt.Sprintf("@%s %s", RightPadTrim(desc, int(float64(listWidth)*0.6)), status)
}
func (stack Stack) FilterValue() string { return stack.Name }

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

type StackSyncResult struct {
	Stack        Stack
	TemplateBody string
	Outtasync    bool
	Err          error
}
