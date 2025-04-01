package ui

import (
	"github.com/dhth/outtasync/internal/types"
)

type (
	DiffFinishedMsg      struct{}
	ShowErrorFinishedMsg struct{}
	ShowDiffFinished     struct{ err error }
	ShowFileFinished     struct{ err error }
)

type TemplateFetchedMsg struct {
	index          int
	templateCode   string
	actualTemplate string
	mismatch       bool
	throttled      bool
	err            error
}

type DriftCheckUpdated struct {
	index     int
	result    types.StackDriftCheckResult
	throttled bool
}

type hideHelpMsg struct{}
