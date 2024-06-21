package ui

type DiffFinishedMsg struct{}
type ShowErrorFinishedMsg struct{}
type CredentialsRefreshedMsg struct{ err error }
type ShowDiffFinished struct{ err error }
type ShowFileFinished struct{ err error }

type TemplateFetchedMsg struct {
	index     int
	stack     Stack
	template  string
	outtaSync bool
	err       error
}
type hideHelpMsg struct{}
