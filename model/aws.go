package model

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	tea "github.com/charmbracelet/bubbletea"
)

func getAWSConfigKey(stack Stack) string {
	return stack.AwsProfile + ":" + stack.AwsRegion
}

func getAWSConfig(profile string, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile))
	return cfg, err

}

func GetCFTemplateBody(awsConfig awsConfig, index int, stack Stack) tea.Cmd {
	return func() tea.Msg {

		if awsConfig.err != nil {
			return TemplateFetchedMsg{index, stack, "", false, awsConfig.err}
		}

		svc := cloudformation.NewFromConfig(awsConfig.config)

		templateInput := cloudformation.GetTemplateInput{
			StackName: &stack.Name,
		}
		templOut, err := svc.GetTemplate(context.TODO(), &templateInput)
		if err != nil {
			return TemplateFetchedMsg{index, stack, "", false, err}
		}

		templBody := *templOut.TemplateBody

		localFile, err := os.ReadFile(stack.Local)
		if err != nil {
			return TemplateFetchedMsg{index, stack, "", false, err}
		}
		localFileContent := string(localFile)
		outtaSync := localFileContent != templBody
		return TemplateFetchedMsg{index, stack, templBody, outtaSync, err}
	}

}
