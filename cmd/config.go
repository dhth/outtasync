package cmd

import (
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/dhth/outtasync/model"
	"gopkg.in/yaml.v3"
)

type T struct {
	Profiles []struct {
		Name   string `yaml:"name"`
		Stacks []struct {
			Name           string  `yaml:"name"`
			Region         string  `yaml:"region"`
			Local          string  `yaml:"local"`
			Tag            *string `yaml:"tag,omitempty"`
			RefreshCommand *string `yaml:"refreshCommand,omitempty"`
		} `yaml:"stacks"`
	} `yaml:"profiles"`
	GlobalRefreshCommand string `yaml:"globalRefreshCommand"`
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			log.Println("Failure reading config")
			os.Exit(1)
		}
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}

func ReadConfig(configFilePath string, profilesToFetch []string) ([]model.Stack, error) {
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
	var rows []model.Stack
	for _, profile := range t.Profiles {
		if len(profilesToFetch) > 0 && !profilesMap[profile.Name] {
			continue
		}
		for _, stack := range profile.Stacks {
			var refreshCmd string
			if stack.RefreshCommand != nil {
				refreshCmd = *stack.RefreshCommand
			} else {
				refreshCmd = globalRefreshCmd
			}
			rows = append(rows, model.Stack{
				Name:           stack.Name,
				AwsProfile:     profile.Name,
				AwsRegion:      stack.Region,
				Template:       "",
				Local:          expandTilde(stack.Local),
				Tag:            stack.Tag,
				RefreshCommand: refreshCmd,
				FetchStatus:    model.StatusUnfetched,
				OuttaSync:      false,
				Err:            nil,
			})
		}
	}
	return rows, nil

}
