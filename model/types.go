package model

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

type Stack struct {
	Name           string
	AwsProfile     string
	AwsRegion      string
	Template       string
	Local          string
	Tag            *string
	RefreshCommand string
	FetchStatus    fetchStatus
	Drifted        bool
	Err            error
}

func (stack Stack) Title() string {
	return fmt.Sprintf("%s", RightPadTrim(stack.Name, listPadding-20))
}
func (stack Stack) Description() string {
	var status string
	switch stack.FetchStatus {
	case StatusFetched:
		switch stack.Drifted {
		case true:
			status = driftedStyle.Render("drifted")
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
	return fmt.Sprintf("@%s %s", RightPadTrim(desc, listPadding-20), status)
}
func (stack Stack) FilterValue() string { return stack.Name }

type delegateKeyMap struct {
	choose             key.Binding
	chooseAll          key.Binding
	refreshCredentials key.Binding
	showDiff           key.Binding
}
