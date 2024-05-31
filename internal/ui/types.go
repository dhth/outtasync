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

const (
	tagWidth = 20
)

type Stack struct {
	Name           string
	AwsProfile     string
	AwsRegion      string
	Template       string
	Local          string
	Tags           []string
	RefreshCommand string
	FetchStatus    fetchStatus
	OuttaSync      bool
	Err            error
}

func (stack Stack) key() string {
	return fmt.Sprintf("%s:%s:%s", stack.AwsProfile, stack.AwsRegion, stack.Name)
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
	var desc string

	var descBudget = int(float64(listWidth) * 0.5)

	if stack.Err != nil {
		desc = RightPadTrim(stack.Err.Error(), int(float64(listWidth)*0.5))
	} else {
		if stack.Tags != nil {
			var tagsLength int
			for _, tag := range stack.Tags {
				nextTag := RightPadTrim(fmt.Sprintf("#%s ", tag), tagWidth)
				if tagsLength+tagWidth > descBudget {
					break
				}
				desc += tagStyle(tag).Render(nextTag)
				tagsLength += tagWidth + 1 // +1 is due to PaddingRight in the style
			}
			for i := 0; i < descBudget-tagsLength; i++ {
				desc += " "
			}
		} else {
			desc = RightPadTrim("@"+stack.AwsProfile, int(float64(listWidth)*0.5))
		}
	}

	return fmt.Sprintf("%s %s",
		desc,
		status,
	)
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
