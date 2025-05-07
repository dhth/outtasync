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

func readConfig(homeDir string, configBytes []byte, stackNameRegex, tagRegex *regexp.Regexp) (types.Config, error) {
	//nolint:prealloc
	var zero types.Config
	cfg := types.OuttasyncConfig{}

	err := yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return zero, err
	}

	config, errorsMsgs := types.ParseConfig(cfg, homeDir, stackNameRegex, tagRegex)

	if len(errorsMsgs) > 0 {
		return zero, fmt.Errorf("%w:\n%s", errInvalidConfig, strings.Join(errorsMsgs, "\n"))
	}

	return config, nil
}
