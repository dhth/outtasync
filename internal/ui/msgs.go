package ui

type (
	DiffFinishedMsg         struct{}
	ShowErrorFinishedMsg    struct{}
	CredentialsRefreshedMsg struct{ err error }
	ShowDiffFinished        struct{ err error }
	ShowFileFinished        struct{ err error }
)

type TemplateFetchedMsg struct {
	index     int
	stackItem stackItem
	template  string
	outtaSync bool
	err       error
}
type hideHelpMsg struct{}
