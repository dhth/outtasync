package cmd

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"

	"flag"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dhth/outtasync/ui"
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var (
	mode         = flag.String("mode", "tui", "the mode to use; possible values: tui/cli")
	pattern      = flag.String("p", "", "regex pattern to filter stack names")
	checkOnStart = flag.Bool("s", false, "whether to check status for all stacks on startup")
)

func Execute() {
	currentUser, err := user.Current()
	var defaultConfigFilePath string
	if err == nil {
		defaultConfigFilePath = fmt.Sprintf("%s/.config/outtasync.yml", currentUser.HomeDir)
	}
	configFilePath := flag.String("config-file", defaultConfigFilePath, "path of the config file")
	profiles := flag.String("profiles", "", "comma separated string of profiles to filter for")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *mode == "" {
		die("mode cannot be empty")
	}

	if *configFilePath == "" {
		die("config-file cannot be empty")
	}

	var regexPattern *regexp.Regexp

	if *pattern != "" {
		regexPattern, err = regexp.Compile(*pattern)
		if err != nil {
			die("Incorrect regex pattern provided: %q\n", err)
		}
	}

	_, err = os.Stat(*configFilePath)
	if os.IsNotExist(err) {
		die(cfgErrSuggestion(fmt.Sprintf("Error: file doesn't exist at %q", *configFilePath)))
	}

	var profilesToFetch []string
	if *profiles != "" {
		profilesToFetch = strings.Split(*profiles, ",")
	}

	stacks, err := ReadConfig(*configFilePath, profilesToFetch, regexPattern)
	if err != nil {
		die(cfgErrSuggestion(fmt.Sprintf("Error reading config: %v", *configFilePath)))
	}
	if len(stacks) == 0 {
		die("No stacks found for the requested parameters")
	}

	awsCfgs := make(map[string]ui.AwsConfig)
	cfClients := make(map[string]ui.AwsCFClient)

	seen := make(map[string]bool)
	for _, stack := range stacks {
		configKey := ui.GetAWSConfigKey(stack)
		if !seen[configKey] {
			cfg, err := ui.GetAWSConfig(stack.AwsProfile, stack.AwsRegion)
			awsCfgs[configKey] = ui.AwsConfig{Config: cfg, Err: err}
			seen[configKey] = true
			if err != nil {
				cfClients[configKey] = ui.AwsCFClient{Err: err}
			} else {
				cfClients[configKey] = ui.AwsCFClient{Client: cloudformation.NewFromConfig(cfg)}
			}
		}
	}

	switch *mode {
	case "tui":
		ui.RenderUI(stacks, awsCfgs, *checkOnStart)
	case "cli":
		ui.ShowResults(stacks, awsCfgs)
	}

}
