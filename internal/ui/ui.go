package ui

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

var errCouldntSetUpDebugLogging = errors.New("couldn't set up debug logging")

func RenderUI(stacks []types.Stack, awsCfgs map[string]aws.Config, checkOnStart bool) error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntSetUpDebugLogging, err.Error())
		}
		defer f.Close()
	}

	p := tea.NewProgram(InitialModel(stacks, awsCfgs, checkOnStart), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}
