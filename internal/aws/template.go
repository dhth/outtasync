package aws

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/dhth/outtasync/internal/types"
	"github.com/dhth/outtasync/internal/utils"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

const stackSyncCheckTimeoutSecs = 5

func CompareStackTemplateCode(
	cfClient CFClient,
	stackName, stackKey, templatePath string,
	remoteCallHeaders []types.RemoteCallHeaders,
	computeDiff bool,
) types.TemplateCheckResult {
	if cfClient.Err != nil {
		return types.TemplateCheckResult{
			StackKey: stackKey,
			Err:      cfClient.Err,
		}
	}

	client := cfClient.Client

	templateInput := cf.GetTemplateInput{
		StackName: &stackName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), stackSyncCheckTimeoutSecs*time.Second)
	defer cancel()

	templOut, err := client.GetTemplate(ctx, &templateInput)
	if err != nil {
		return types.TemplateCheckResult{
			StackKey: stackKey,
			Err:      err,
		}
	}

	actualTemplate := *templOut.TemplateBody

	// to deal with "No newline at end of file" issues while diffing
	if !strings.HasSuffix(actualTemplate, "\n") {
		actualTemplate = actualTemplate + "\n"
	}

	var templateCode string

	if strings.HasPrefix(templatePath, "https://") {
		var headers [][2]string
		for _, header := range remoteCallHeaders {
			headers = append(headers, [2]string{os.ExpandEnv(header.Key), os.ExpandEnv(header.Value)})
		}

		respBytes, err := utils.GetHTTPResponse(templatePath, headers)
		if err != nil {
			return types.TemplateCheckResult{
				StackKey: stackKey,
				Err:      err,
			}
		}
		templateCode = string(respBytes)
	} else {
		localFile, err := os.ReadFile(templatePath)
		if err != nil {
			return types.TemplateCheckResult{
				StackKey: stackKey,
				Err:      err,
			}
		}
		templateCode = string(localFile)
	}

	// to deal with "No newline at end of file" issues while diffing
	if !strings.HasSuffix(templateCode, "\n") {
		templateCode = templateCode + "\n"
	}

	mismatch := templateCode != actualTemplate

	var diff []byte
	var diffErr error
	if computeDiff {
		diff, diffErr = utils.GetDiff(templateCode, actualTemplate)
	}

	return types.TemplateCheckResult{
		StackKey:       stackKey,
		TemplateCode:   templateCode,
		ActualTemplate: actualTemplate,
		Diff:           diff,
		DiffErr:        diffErr,
		Mismatch:       mismatch,
	}
}
