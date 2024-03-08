package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"flag"

	"github.com/dhth/outtasync/ui"
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Execute() {
	currentUser, err := user.Current()
	var defaultConfigFilePath string
	if err == nil {
		defaultConfigFilePath = fmt.Sprintf("%s/.config/outtasync.yml", currentUser.HomeDir)
	}
	configFilePath := flag.String("config-file", defaultConfigFilePath, "path of the config file")
	profiles := flag.String("profiles", "", "comma separated string of profiles to filter for")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *configFilePath == "" {
		die("config-file cannot be empty")
	}

	_, err = os.Stat(*configFilePath)
	if os.IsNotExist(err) {
		die(cfgErrSuggestion(fmt.Sprintf("Error: file doesn't exist at %q", *configFilePath)))
	}

	var profilesToFetch []string
	if *profiles != "" {
		profilesToFetch = strings.Split(*profiles, ",")
	}

	stacks, err := ReadConfig(*configFilePath, profilesToFetch)
	if err != nil {
		die(cfgErrSuggestion(fmt.Sprintf("Error reading config: %v", *configFilePath)))
	}
	if len(stacks) == 0 {
		die(cfgErrSuggestion(fmt.Sprintf("No stacks found for the requested parameters")))
	}
	ui.RenderUI(stacks)

}
