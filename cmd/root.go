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
	profiles := flag.String("profiles", "", "profiles to filter for")

	flag.Parse()

	_, err = os.Stat(*configFilePath)
	if os.IsNotExist(err) {
		die("Error: file not found at %s\n", *configFilePath)
	}
	var profilesToFetch []string
	if *profiles != "" {
		profilesToFetch = strings.Split(*profiles, ",")
	}

	stacks, err := ReadConfig(*configFilePath, profilesToFetch)
	if err != nil {
		die("Error reading config: %v\n", err)
	}
	if len(stacks) == 0 {
		die("No stacks found for the requested parameters")
	}
	ui.RenderUI(stacks)

}
