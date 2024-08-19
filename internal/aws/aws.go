package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dhth/outtasync/internal/types"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type Config struct {
	Config aws.Config
	Err    error
}

type CFClient struct {
	Client *cloudformation.Client
	Err    error
}

func GetAWSConfig(profile string, region string) (aws.Config, error) {
	if profile == "default" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region))
		return cfg, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile))
	return cfg, err
}

func CheckStackSyncStatus(awsConfig Config, stack types.Stack) types.StackSyncResult {
	if awsConfig.Err != nil {
		return types.StackSyncResult{
			Stack: stack,
			Err:   awsConfig.Err,
		}
	}

	svc := cloudformation.NewFromConfig(awsConfig.Config)

	templateInput := cloudformation.GetTemplateInput{
		StackName: &stack.Name,
	}
	templOut, err := svc.GetTemplate(context.TODO(), &templateInput)
	if err != nil {
		return types.StackSyncResult{
			Stack: stack,
			Err:   err,
		}
	}

	templBody := *templOut.TemplateBody

	localFile, err := os.ReadFile(stack.Local)
	if err != nil {
		return types.StackSyncResult{
			Stack: stack,
			Err:   err,
		}
	}
	localFileContent := string(localFile)
	outtaSync := localFileContent != templBody

	return types.StackSyncResult{
		Stack:        stack,
		TemplateBody: templBody,
		Outtasync:    outtaSync,
	}
}
