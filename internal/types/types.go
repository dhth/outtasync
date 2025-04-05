package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dhth/outtasync/internal/utils"
	"gopkg.in/yaml.v3"
)

const (
	cfgSrcSharedProfilePrefix = "profile:"
	cfgSrcAssumeRolePrefix    = "assume:"
)

var (
	errConfigSourceEmpty     = errors.New("config source is empty")
	errIncorrectConfigSource = errors.New("incorrect config source provided")
)

type ConfigSourceKind uint

const (
	Env ConfigSourceKind = iota
	SharedProfile
	AssumeRole
)

type ConfigSource struct {
	Kind  ConfigSourceKind
	Value string
}

func (cs ConfigSource) Display() string {
	var value string
	switch cs.Kind {
	case Env:
		value = "env"
	case SharedProfile:
		value = fmt.Sprintf("profile:%s", cs.Value)
	case AssumeRole:
		value = fmt.Sprintf("assume:%s", cs.Value)
	}

	return value
}

type Config struct {
	Stacks []StackConfig `yaml:"stacks"`
}

type StackConfig struct {
	Name                      string                      `yaml:"name"`
	ConfigSource              string                      `yaml:"config_source"`
	Arn                       string                      `yaml:"arn"`
	TemplatePath              *string                     `yaml:"template_path,omitempty"`
	TemplateRemoteCallHeaders []TemplateRemoteCallHeaders `yaml:"template_remote_call_headers,omitempty"`
	Tags                      []string                    `yaml:"tags,omitempty"`
}

type TemplateRemoteCallHeaders struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func (sc StackConfig) getConfigSource() (ConfigSource, error) {
	return ParseConfigSource(sc.ConfigSource)
}

type Stack struct {
	Name                      string
	Arn                       string
	ConfigSource              ConfigSource
	TemplatePath              *string
	TemplateRemoteCallHeaders []TemplateRemoteCallHeaders
	Tags                      []string
}

func (s Stack) GetConfigRepr() StackConfig {
	return StackConfig{
		Name:         s.Name,
		ConfigSource: s.ConfigSource.Display(),
		Arn:          s.Arn,
		Tags:         s.Tags,
	}
}

func (s Stack) Key() string {
	return s.Arn
}

func (s Stack) AWSConfigKey() string {
	return s.ConfigSource.Value
}

type TemplateCheckResult struct {
	StackKey       string
	TemplateCode   string
	ActualTemplate string
	Diff           []byte
	DiffErr        error
	Mismatch       bool
	Err            error
}

type StackDriftCheckResult struct {
	Stack  Stack
	Output *cf.DescribeStackDriftDetectionStatusOutput
	Err    error
}

func ParseConfigSource(value string) (ConfigSource, error) {
	var zero ConfigSource
	if strings.TrimSpace(value) == "" {
		return zero, errConfigSourceEmpty
	}

	if value == "env" {
		return ConfigSource{Env, "env"}, nil
	}

	if strings.HasPrefix(value, cfgSrcSharedProfilePrefix) {
		value := strings.TrimPrefix(value, cfgSrcSharedProfilePrefix)
		if strings.TrimSpace(value) == "" {
			return zero, errConfigSourceEmpty
		}
		return ConfigSource{
			SharedProfile,
			value,
		}, nil
	}

	if strings.HasPrefix(value, cfgSrcAssumeRolePrefix) {
		value := strings.TrimPrefix(value, cfgSrcAssumeRolePrefix)
		if strings.TrimSpace(value) == "" {
			return zero, errConfigSourceEmpty
		}
		return ConfigSource{
			SharedProfile,
			value,
		}, nil
	}

	return zero, errIncorrectConfigSource
}

func EncodeConfig(stacks []Stack) ([]byte, error) {
	stackConfigs := make([]StackConfig, len(stacks))
	for i, st := range stacks {
		stackConfigs[i] = st.GetConfigRepr()
	}

	config := Config{
		Stacks: stackConfigs,
	}

	var zero []byte
	configBytes := bytes.Buffer{}

	enc := yaml.NewEncoder(&configBytes)
	enc.SetIndent(2)
	err := enc.Encode(&config)
	if err != nil {
		return zero, err
	}

	return configBytes.Bytes(), nil
}

func ParseStackConfig(config StackConfig, homeDir string) (Stack, []error) {
	var errors []error
	var zero Stack

	configSource, err := config.getConfigSource()
	if err != nil {
		errors = append(errors, err)
	}

	var local *string
	if config.TemplatePath != nil {
		l := utils.ExpandTilde(*config.TemplatePath, homeDir)
		local = &l
	}

	if len(errors) > 0 {
		return zero, errors
	}

	return Stack{
		Name:                      config.Name,
		Arn:                       config.Arn,
		ConfigSource:              configSource,
		TemplatePath:              local,
		TemplateRemoteCallHeaders: config.TemplateRemoteCallHeaders,
		Tags:                      config.Tags,
	}, nil
}

type CheckOutputFormat uint8

const (
	Default CheckOutputFormat = iota
	Delimited
	HTML
)

type CheckHTMLOutputConfig struct {
	Title    string
	Template string
	Open     bool
}

func ParseCheckOutputFormat(value string) (CheckOutputFormat, bool) {
	var zero CheckOutputFormat
	switch value {
	case "default":
		return Default, true
	case "delimited":
		return Delimited, true
	case "html":
		return HTML, true
	default:
		return zero, false
	}
}

func CheckOutputFormatPossibleValues() []string {
	return []string{
		"default",
		"delimited",
		"html",
	}
}
