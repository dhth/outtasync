package ui

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func GetAWSConfigKey(stack Stack) string {
	return stack.AwsProfile + ":" + stack.AwsRegion
}

func GetAWSConfig(profile string, region string) (aws.Config, error) {
	if profile == "default" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region))
		return cfg, err
	} else {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithSharedConfigProfile(profile))
		return cfg, err
	}
}

func CheckStackSyncStatus(awsConfig AwsConfig, stack Stack) StackSyncResult {
	if awsConfig.Err != nil {
		return StackSyncResult{stack, "", false, awsConfig.Err}
	}

	svc := cloudformation.NewFromConfig(awsConfig.Config)

	templateInput := cloudformation.GetTemplateInput{
		StackName: &stack.Name,
	}
	templOut, err := svc.GetTemplate(context.TODO(), &templateInput)
	if err != nil {
		return StackSyncResult{stack, "", false, err}
	}

	templBody := *templOut.TemplateBody

	localFile, err := os.ReadFile(stack.Local)
	if err != nil {
		return StackSyncResult{stack, "", false, err}
	}
	localFileContent := string(localFile)
	outtaSync := localFileContent != templBody
	return StackSyncResult{stack, templBody, outtaSync, err}
}
