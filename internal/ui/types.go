package ui

import (
	"fmt"

	"github.com/dhth/outtasync/internal/types"
)

type fetchStatus uint

const (
	statusUnfetched fetchStatus = iota
	statusFetching
	statusFetched
	statusFailure
)

const (
	tagWidth         = 20
	stackNamePadding = 80
)

type stackItem struct {
	stack       types.Stack
	fetchStatus fetchStatus
	outtaSync   bool
	err         error
}

func (si stackItem) Title() string {
	var status string
	switch si.fetchStatus {
	case statusFetched:
		switch si.outtaSync {
		case true:
			status = outtaSyncStyle.Render("outta sync")
		case false:
			status = insSyncStyle.Render("in sync")
		}
	case statusFetching:
		status = fetchingStyle.Render("...")
	case statusFailure:
		status = errorStyle.Render("error")
	}
	name := RightPadTrim(si.stack.Name, stackNamePadding)
	return fmt.Sprintf("%s %s",
		name,
		status,
	)
}

func (si stackItem) Description() string {
	var desc string

	if si.err != nil {
		return si.err.Error()
	}

	if si.stack.Tags != nil {
		for _, tag := range si.stack.Tags {
			nextTag := RightPadTrim(fmt.Sprintf("#%s ", tag), tagWidth)
			desc += tagStyle(tag).Render(nextTag)
		}
		return desc
	}

	return fmt.Sprintf("@%s", si.stack.AwsProfile)
}

func (si stackItem) FilterValue() string { return si.stack.Name }
