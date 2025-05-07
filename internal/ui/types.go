package ui

import (
	"fmt"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/dhth/outtasync/internal/types"
	"github.com/dhth/outtasync/internal/utils"
)

type pane uint

const (
	stacksList pane = iota
	codeMismatchStacksList
	driftedStacksList
	erroredStacksList
	errorDetailsPane
	helpPane
)

type syncCheckStatus uint

const (
	syncStatusNotChecked syncCheckStatus = iota
	syncStatusInProgress
	syncStatusChecked
	syncStatusCheckFailed
)

type driftCheckStatus uint

const (
	driftNotChecked driftCheckStatus = iota
	driftCheckInProgress
	driftChecked
)

const (
	tagWidth         = 20
	stackNamePadding = 80
)

type stackItem struct {
	stack            types.Stack
	templateCode     string
	actualTemplate   string
	syncCheckStatus  syncCheckStatus
	outtaSync        bool
	syncErr          error
	driftCheckStatus driftCheckStatus
	driftOutput      *cf.DescribeStackDriftDetectionStatusOutput
	driftErr         error
}

func (si *stackItem) hasDrifted() bool {
	if si.driftOutput == nil {
		return false
	}

	return si.driftOutput.StackDriftStatus == cftypes.StackDriftStatusDrifted
}

func (si stackItem) Title() string {
	var status string
	switch si.syncCheckStatus {
	case syncStatusChecked:
		switch si.outtaSync {
		case true:
			status = outtaSyncStyle.Render("outtasync")
		case false:
			status = insSyncStyle.Render("✓")
		}
	case syncStatusInProgress:
		status = fetchingStyle.Render("...")
	case syncStatusCheckFailed:
		status = errorStyle.Render("error")
	}

	var driftStatus string
	if si.driftErr != nil {
		driftStatus = errorStyle.Render("error")
	} else {
		switch si.driftCheckStatus {
		case driftCheckInProgress:
			driftStatus = driftCheckInProgressStyle.Render("...")
		case driftChecked:
			if si.driftOutput != nil {
				ds := fmt.Sprintf("%v", si.driftOutput.StackDriftStatus)
				var rc string
				switch si.driftOutput.StackDriftStatus {
				case cftypes.StackDriftStatusDrifted:
					ds = driftedStyle.Render("drifted")
					if si.driftOutput.DriftedStackResourceCount != nil {
						if *si.driftOutput.DriftedStackResourceCount == 1 {
							rc = " (1 resource) "
						} else {
							rc = fmt.Sprintf(" (%d resources) ", *si.driftOutput.DriftedStackResourceCount)
						}
					}
				case cftypes.StackDriftStatusInSync:
					ds = notDriftedStyle.Render("✓")
				case cftypes.StackDriftStatusUnknown:
					ds = unknownDriftStatusStyle.Render("unknown")
				}

				var dr string
				if si.driftOutput.DetectionStatusReason != nil {
					dr = driftReasonStyle.Render(*si.driftOutput.DetectionStatusReason)
				}

				driftStatus = fmt.Sprintf("%s%s %s", ds, rc, dr)
			}
		}
	}

	name := utils.RightPadTrim(si.stack.Name, stackNamePadding)
	if status == "" {
		status = utils.RightPadTrim("", 11)
	}
	if driftStatus != "" {
		driftStatus = "  " + driftStatus
	}

	return fmt.Sprintf("%s %s%s",
		name,
		status,
		driftStatus,
	)
}

func (si stackItem) Description() string {
	if si.stack.Tags == nil {
		return ""
	}

	var desc string
	for _, tag := range si.stack.Tags {
		nextTag := utils.RightPadTrim(fmt.Sprintf("#%s ", tag), tagWidth)
		desc += tagStyle(tag).Render(nextTag)
	}
	return desc
}

func (si stackItem) FilterValue() string { return si.stack.Name }
