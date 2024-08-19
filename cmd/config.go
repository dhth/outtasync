package cmd

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dhth/outtasync/internal/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Profiles []struct {
		Name   string `yaml:"name"`
		Stacks []struct {
			Name           string   `yaml:"name"`
			Region         string   `yaml:"region"`
			Local          string   `yaml:"local"`
			Tags           []string `yaml:"tags,omitempty"`
			RefreshCommand *string  `yaml:"refreshCommand,omitempty"`
		} `yaml:"stacks"`
	} `yaml:"profiles"`
	GlobalRefreshCommand string `yaml:"globalRefreshCommand"`
}

func expandTilde(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func readConfig(homeDir string, configBytes []byte, profilesToFetch []string, tagsToFetch []string, pattern *regexp.Regexp,
) ([]types.Stack, error) {
	var stacks []types.Stack
	cfg := Config{}

	err := yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return stacks, err
	}

	profilesMap := make(map[string]bool)
	for _, p := range profilesToFetch {
		profilesMap[p] = true
	}

	globalRefreshCmd := cfg.GlobalRefreshCommand
	for _, profile := range cfg.Profiles {

		if len(profilesToFetch) > 0 && !profilesMap[profile.Name] {
			continue
		}

		for _, stack := range profile.Stacks {
			if pattern != nil && !pattern.MatchString(stack.Name) {
				continue
			}

			if len(tagsToFetch) > 0 {
				if stack.Tags == nil {
					continue
				}

				stackTagsMap := make(map[string]bool)
				for _, tag := range stack.Tags {
					stackTagsMap[tag] = true
				}

				tagNotInStack := false
				for _, tagToFetch := range tagsToFetch {
					if !stackTagsMap[tagToFetch] {
						tagNotInStack = true
						break
					}
				}
				if tagNotInStack {
					continue
				}

			}
			var refreshCmd string
			if stack.RefreshCommand != nil {
				refreshCmd = *stack.RefreshCommand
			} else {
				refreshCmd = globalRefreshCmd
			}

			stacks = append(stacks, types.Stack{
				Name:           stack.Name,
				AwsProfile:     profile.Name,
				AwsRegion:      stack.Region,
				Local:          expandTilde(stack.Local, homeDir),
				Tags:           stack.Tags,
				RefreshCommand: refreshCmd,
			})
		}
	}
	return stacks, nil
}
