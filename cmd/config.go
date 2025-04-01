package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dhth/outtasync/internal/types"
	"gopkg.in/yaml.v3"
)

var errInvalidConfig = errors.New("invalid config provided")

func readConfig(homeDir string, configBytes []byte, stackNameRegex, tagRegex *regexp.Regexp) ([]types.Stack, error) {
	//nolint:prealloc
	var stacks []types.Stack
	cfg := types.Config{}

	err := yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return stacks, err
	}

	var errors []string
	for i, sc := range cfg.Stacks {
		stack, errs := types.ParseStackConfig(sc, homeDir)
		if len(errs) > 0 {
			errors = append(errors, fmt.Sprintf("- invalid config for stack at index %d: %v", i, errs))
			continue
		}

		if stackNameRegex != nil && !stackNameRegex.Match([]byte(stack.Name)) {
			continue
		}

		if tagRegex != nil {
			tagMatch := false
			for _, tag := range stack.Tags {
				if tagRegex.Match([]byte(tag)) {
					tagMatch = true
					break
				}
			}
			if !tagMatch {
				continue
			}
		}

		stacks = append(stacks, stack)
	}

	if len(errors) > 0 {
		return stacks, fmt.Errorf("%w:\n%s", errInvalidConfig, strings.Join(errors, "\n"))
	}

	return stacks, nil
}
