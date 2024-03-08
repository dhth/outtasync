package cmd

import "fmt"

var (
	configSampleFormat = `globalRefreshCommand: aws sso login --sso-session sessionname
profiles:
- name: qa
  stacks:
  - name: bingo-service-qa
    local: ~/projects/bingo-service/cloudformation/infrastructure.yml
    region: eu-central-1
    refreshCommand: aws sso login --profile qa1
  - name: papaya-service-qa
    local: ~/projects/papaya-service/cloudformation/service.yml
    region: eu-central-1
  - name: racoon-service-qa
    local: ~/projects/racoon-service/cloudformation/service.yml
    region: eu-central-1
- name: prod
  stacks:
  - name: brb-dll-prod
    local: ~/projects/brd-dll-service/cloudformation/service.yml
    region: eu-central-1
    refreshCommand: aws sso login --profile rgb-prod
  - name: galactus-service-prod
    local: ~/projects/galactus-service/cloudformation/service.yml
    region: eu-central-1`
	helpText = `Identify cloudformation stacks that have gone out of sync with the state represented by their stack files

Usage: outtasync [flags]`
)

func cfgErrSuggestion(msg string) string {
	return fmt.Sprintf(`%s

Make sure to structure the yml config file as follows:

%s

Use "outtasync -help" for more information`,
		msg,
		configSampleFormat,
	)
}
