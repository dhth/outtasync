package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/ui"
)

const (
	configFileName = "outtasync/outtasync.yml"
	helpText       = `Identify cloudformation stacks that have gone out of sync with the state represented by their stack files.

Usage: outtasync [flags]`
)

var (
	mode         = flag.String("mode", "tui", "the mode to use; possible values: tui/cli")
	pattern      = flag.String("p", "", "regex pattern to filter stack names")
	profiles     = flag.String("profiles", "", "comma separated string of profiles to filter for")
	tags         = flag.String("t", "", "comma separated string of tags to filter for; will match stacks that contain all tags specified here")
	checkOnStart = flag.Bool("c", false, "whether to check status for all stacks on startup")
)

var (
	errModeFlagEmpty          = errors.New("mode flag cannot be empty")
	errConfigFileFlagEmpty    = errors.New("config file flag cannot be empty")
	errCouldntGetHomeDir      = errors.New("couldn't get your home directory")
	errCouldntGetConfigDir    = errors.New("couldn't get your default config directory")
	errConfigFileExtIncorrect = errors.New("config file must be a YAML file")
	errConfigFileDoesntExist  = errors.New("config file does not exist")
	errCouldntReadConfigFile  = errors.New("couldn't read config file")
	errCouldntParseConfigFile = errors.New("couldn't parse config file")
	errIncorrectRegexProvided = errors.New("incorrect regex provided")
	errNoStacksFound          = errors.New("no stacks found")
)

func Execute() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGetHomeDir, err.Error())
	}

	defaultConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGetConfigDir, err.Error())
	}
	defaultConfigFilePath := filepath.Join(defaultConfigDir, configFileName)
	configFilePath := flag.String("config-file", defaultConfigFilePath, "path of the config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *mode == "" {
		return fmt.Errorf("%w", errModeFlagEmpty)
	}

	if *configFilePath == "" {
		return fmt.Errorf("%w", errConfigFileFlagEmpty)
	}

	configPathFull := expandTilde(*configFilePath, userHomeDir)

	var regexPattern *regexp.Regexp

	if *pattern != "" {
		regexPattern, err = regexp.Compile(*pattern)
		if err != nil {
			return fmt.Errorf("%w: %s", errIncorrectRegexProvided, err.Error())
		}
	}

	if filepath.Ext(configPathFull) != ".yml" && filepath.Ext(configPathFull) != ".yaml" {
		return errConfigFileExtIncorrect
	}

	_, err = os.Stat(configPathFull)
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", errConfigFileDoesntExist, err.Error())
	}

	var profilesToFetch []string
	if *profiles != "" {
		profilesToFetch = strings.Split(*profiles, ",")
	}

	var tagsToFetch []string
	if *tags != "" {
		tagsToFetch = strings.Split(*tags, ",")
	}

	configBytes, err := os.ReadFile(configPathFull)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntReadConfigFile, err.Error())
	}

	stacks, err := readConfig(userHomeDir, configBytes, profilesToFetch, tagsToFetch, regexPattern)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntParseConfigFile, err.Error())
	}

	if len(stacks) == 0 {
		return fmt.Errorf("%w", errNoStacksFound)
	}

	awsCfgs := make(map[string]aws.Config)
	cfClients := make(map[string]aws.CFClient)

	seen := make(map[string]bool)
	for _, stack := range stacks {
		configKey := stack.AWSConfigKey()
		if !seen[configKey] {
			cfg, err := aws.GetAWSConfig(stack.AwsProfile, stack.AwsRegion)
			awsCfgs[configKey] = aws.Config{Config: cfg, Err: err}
			seen[configKey] = true
			if err != nil {
				cfClients[configKey] = aws.CFClient{Err: err}
			} else {
				cfClients[configKey] = aws.CFClient{Client: cloudformation.NewFromConfig(cfg)}
			}
		}
	}

	switch *mode {
	case "tui":
		err = ui.RenderUI(stacks, awsCfgs, *checkOnStart)
		if err != nil {
			return err
		}
	case "cli":
		showResults(stacks, awsCfgs)
	}

	return nil
}
