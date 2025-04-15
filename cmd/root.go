package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dhth/outtasync/internal/aws"
	cli "github.com/dhth/outtasync/internal/cli"
	"github.com/dhth/outtasync/internal/types"
	"github.com/dhth/outtasync/internal/ui"
	"github.com/dhth/outtasync/internal/utils"
	"github.com/spf13/cobra"
)

const (
	configFileName = "outtasync/outtasync.yml"
)

var (
	errCouldntGetUserHomeDir   = errors.New("couldn't get your home directory")
	errCouldntGetUserConfigDir = errors.New("couldn't get your config directory")
	ErrCouldntReadConfigFile   = errors.New("couldn't read config file")
	errIncorrectFormatProvided = errors.New("incorrect format provided")
	errNothingToCheck          = errors.New("nothing to check")
	errTemplateFileNotHTML     = errors.New("template is not an HTML file")
	errCouldntReadTemplateFile = errors.New("couldn't read template file")
)

func Execute() error {
	rootCmd, err := NewRootCommand()
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func NewRootCommand() (*cobra.Command, error) {
	var (
		configPath     string
		configPathFull string
		homeDir        string

		nameFilterStr               string
		tagsFilterStr               string
		checkShowProgressIndicator  bool
		compareWithCode             bool
		checkDrift                  bool
		checkOutputFormat           string
		listNegativesOnly           bool
		checkHTMLOutputTitle        string
		checkHTMLOutputTemplateFile string
		checkHTMLOpen               bool

		genConfigConfigSource string
		genConfigFilterStr    string
		genConfigTags         string
	)

	rootCmd := &cobra.Command{
		Use: "outtasync",
		Short: `outtasync helps you identify Cloudformation stacks that have drifted or gone out of sync
with the state represented by their template files.`,
		SilenceUsage: true,
	}

	checkCmd := &cobra.Command{
		Use:          "check",
		Short:        "Check sync and drift status for stacks",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			var stackNameRegex *regexp.Regexp
			var err error
			var tagRegex *regexp.Regexp

			if nameFilterStr != "" {
				stackNameRegex, err = regexp.Compile(nameFilterStr)
				if err != nil {
					return err
				}
			}
			if tagsFilterStr != "" {
				tagRegex, err = regexp.Compile(tagsFilterStr)
				if err != nil {
					return err
				}
			}

			outputFormat, ok := types.ParseCheckOutputFormat(checkOutputFormat)
			if !ok {
				return errIncorrectFormatProvided
			}

			if !compareWithCode && !checkDrift {
				return fmt.Errorf("%w; either provide --compare-template or --check-drift or both", errNothingToCheck)
			}

			configPathFull = utils.ExpandTilde(configPath, homeDir)
			configBytes, err := os.ReadFile(configPathFull)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
			}

			config, err := readConfig(homeDir, configBytes, stackNameRegex, tagRegex)
			if err != nil {
				return err
			}

			if len(config.Stacks) == 0 {
				return errNothingToCheck
			}

			cfClients := make(map[string]aws.CFClient)

			seen := make(map[string]bool)
			for _, stack := range config.Stacks {
				configKey := stack.AWSConfigKey()
				if seen[configKey] {
					continue
				}

				cfg, err := aws.GetAWSConfig(stack.ConfigSource)
				seen[configKey] = true
				if err != nil {
					cfClients[configKey] = aws.CFClient{Err: err}
				} else {
					cfClients[configKey] = aws.CFClient{Client: cf.NewFromConfig(cfg)}
				}
			}

			var htmlOutputConf types.CheckHTMLOutputConfig
			if outputFormat == types.HTML {
				var htmlOutputTemplate string
				if checkHTMLOutputTemplateFile != "" {
					if !strings.HasSuffix(checkHTMLOutputTemplateFile, ".html") {
						return errTemplateFileNotHTML
					}

					htmlTemplateFileFull := utils.ExpandTilde(checkHTMLOutputTemplateFile, homeDir)
					templateBytes, err := os.ReadFile(htmlTemplateFileFull)
					if err != nil {
						return fmt.Errorf("%w: %s", errCouldntReadTemplateFile, err.Error())
					}
					htmlOutputTemplate = string(templateBytes)
				}
				htmlOutputConf = types.CheckHTMLOutputConfig{
					Title:    checkHTMLOutputTitle,
					Template: htmlOutputTemplate,
					Open:     checkHTMLOpen,
				}
			}
			return cli.ShowCheckResults(config,
				cfClients,
				compareWithCode,
				checkDrift,
				checkShowProgressIndicator,
				outputFormat,
				listNegativesOnly,
				&htmlOutputConf,
			)
		},
	}

	tuiCmd := &cobra.Command{
		Use:          "tui",
		Short:        "Open outtasync's tui",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			var stackNameRegex *regexp.Regexp
			var err error
			var tagRegex *regexp.Regexp

			if nameFilterStr != "" {
				stackNameRegex, err = regexp.Compile(nameFilterStr)
				if err != nil {
					return err
				}
			}
			if tagsFilterStr != "" {
				tagRegex, err = regexp.Compile(tagsFilterStr)
				if err != nil {
					return err
				}
			}

			configPathFull = utils.ExpandTilde(configPath, homeDir)
			configBytes, err := os.ReadFile(configPathFull)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCouldntReadConfigFile, err)
			}

			config, err := readConfig(homeDir, configBytes, stackNameRegex, tagRegex)
			if err != nil {
				return err
			}

			if len(config.Stacks) == 0 {
				return nil
			}

			cfClients := make(map[string]aws.CFClient)

			seen := make(map[string]bool)
			for _, stack := range config.Stacks {
				configKey := stack.AWSConfigKey()
				if seen[configKey] {
					continue
				}

				cfg, err := aws.GetAWSConfig(stack.ConfigSource)
				seen[configKey] = true
				if err != nil {
					cfClients[configKey] = aws.CFClient{Err: err}
				} else {
					cfClients[configKey] = aws.CFClient{Client: cf.NewFromConfig(cfg)}
				}
			}

			return ui.RenderUI(config, cfClients)
		},
	}

	configCmd := &cobra.Command{
		Use:          "config",
		Short:        "Work with outtasync's config",
		SilenceUsage: true,
	}

	generateConfigCmd := &cobra.Command{
		Use:          "generate",
		Short:        "generate sample config",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			var filterRegex *regexp.Regexp
			var err error

			if genConfigFilterStr != "" {
				filterRegex, err = regexp.Compile(genConfigFilterStr)
				if err != nil {
					return err
				}
			}
			configSource, err := types.ParseConfigSource(genConfigConfigSource)
			if err != nil {
				return err
			}

			cfg, err := aws.GetAWSConfig(configSource)
			if err != nil {
				return err
			}

			client := cf.NewFromConfig(cfg)

			tags := strings.Split(genConfigTags, ",")

			stacks, err := aws.GetStacksForAccount(client, filterRegex, configSource, tags)
			if err != nil {
				return err
			}

			if len(stacks) == 0 {
				return nil
			}

			configBytes, err := types.EncodeConfig(stacks)
			if err != nil {
				return err
			}

			fmt.Print(string(configBytes))

			return nil
		},
	}

	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserHomeDir, err.Error())
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserConfigDir, err.Error())
	}

	defaultConfigPath := filepath.Join(configDir, configFileName)

	checkCmd.Flags().StringVarP(&configPath, "config-file", "c", defaultConfigPath, "location of outtasync's config file")
	checkCmd.Flags().StringVarP(&nameFilterStr, "name-filter", "n", "", "regex for name(s) (configured in outtasync's config) to filter stacks by")
	checkCmd.Flags().StringVarP(&tagsFilterStr, "tags-filter", "t", "", "regex for tag(s) to filter stacks by")
	checkCmd.Flags().BoolVarP(&checkShowProgressIndicator, "progress-indicator", "p", true, "whether to show progress indicator (only applicable in cli mode)")
	checkCmd.Flags().BoolVarP(&compareWithCode, "compare-template", "T", false, "compare actual template with template code (only applicable in cli mode)")
	checkCmd.Flags().BoolVarP(&checkDrift, "check-drift", "D", true, "check drift status (only applicable in cli mode)")
	checkCmd.Flags().StringVarP(&checkOutputFormat, "format", "f", "default", fmt.Sprintf("output format [possible values: %s]", strings.Join(types.CheckOutputFormatPossibleValues(), ", ")))
	checkCmd.Flags().BoolVarP(&listNegativesOnly, "list-negatives-only", "N", false, "list negatives only")
	checkCmd.Flags().StringVar(&checkHTMLOutputTitle, "html-title", "outtasync", "title of the html output")
	checkCmd.Flags().StringVar(&checkHTMLOutputTemplateFile, "html-template-file", "", "location of the template file to use for html output")
	checkCmd.Flags().BoolVarP(&checkHTMLOpen, "html-open", "o", false, "open html output in browser instead of outputting to stdout")

	tuiCmd.Flags().StringVarP(&configPath, "config-file", "c", defaultConfigPath, "location of outtasync's config file")
	tuiCmd.Flags().StringVarP(&nameFilterStr, "name-filter", "n", "", "regex for name(s) (configured in outtasync's config) to filter stacks by")
	tuiCmd.Flags().StringVarP(&tagsFilterStr, "tags-filter", "t", "", "regex for tag(s) to filter stacks by")

	generateConfigCmd.Flags().StringVarP(&genConfigConfigSource, "config-source", "c", "env", "config source to use")
	generateConfigCmd.Flags().StringVarP(&genConfigFilterStr, "name-filter", "n", "", "regex for name(s) to filter stacks by")
	generateConfigCmd.Flags().StringVarP(&genConfigTags, "tags", "t", "", "comma separated list of tags to use")

	configCmd.AddCommand(generateConfigCmd)

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(tuiCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}
