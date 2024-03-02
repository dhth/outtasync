package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dhth/outtasync/ui"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "outtasync",
		Short:   "A TUI to identify cloudformation stacks that have gone out of sync with the state represented by their stack files.",
		Version: "0.1",
		Run: func(cmd *cobra.Command, args []string) {
			configFilePath, err := cmd.Flags().GetString("config-file")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			_, err = os.Stat(configFilePath)
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: file not found at %s\n", configFilePath)
				os.Exit(1)
			}
			stacks, err := ReadConfig(configFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
				os.Exit(1)
			}
			ui.RenderUI(stacks)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var configFilePath string
	currentUser, err := user.Current()
	var defaultConfigFilePath string
	if err == nil {
		defaultConfigFilePath = fmt.Sprintf("%s/.config/outtasync.yml", currentUser.HomeDir)
	}
	rootCmd.Flags().StringVarP(&configFilePath, "config-file", "c", defaultConfigFilePath, "Path of the config file")
}
