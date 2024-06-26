package cmd

import (
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/dhth/outtasync/internal/ui"
	"gopkg.in/yaml.v3"
)

type T struct {
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

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			os.Exit(1)
		}
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}

func readConfig(configFilePath string,
	profilesToFetch []string,
	tagsToFetch []string,
	pattern *regexp.Regexp) ([]ui.Stack, error) {

	localFile, err := os.ReadFile(expandTilde(configFilePath))
	if err != nil {
		os.Exit(1)
	}
	t := T{}
	err = yaml.Unmarshal(localFile, &t)
	if err != nil {
		return nil, err
	}
	profilesMap := make(map[string]bool)
	for _, p := range profilesToFetch {
		profilesMap[p] = true
	}

	globalRefreshCmd := t.GlobalRefreshCommand
	var rows []ui.Stack
	for _, profile := range t.Profiles {

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
			rows = append(rows, ui.Stack{
				Name:           stack.Name,
				AwsProfile:     profile.Name,
				AwsRegion:      stack.Region,
				Template:       "",
				Local:          expandTilde(stack.Local),
				Tags:           stack.Tags,
				RefreshCommand: refreshCmd,
				FetchStatus:    ui.StatusUnfetched,
				OuttaSync:      false,
				Err:            nil,
			})
		}
	}
	return rows, nil

}
