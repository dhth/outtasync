package model

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	tea "github.com/charmbracelet/bubbletea"
)

func GetCFTemplateBody(index int, stack Stack) tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(stack.AwsRegion),
			config.WithSharedConfigProfile(stack.AwsProfile))
		if err != nil {
			return TemplateFetchedMsg{index, stack, "", false, err}
		}

		svc := cloudformation.NewFromConfig(cfg)

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
