package aws

import (
	"context"
	"regexp"
	"time"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/dhth/outtasync/internal/types"
)

const getStacksMaxPages = 5

func GetStacksForAccount(
	client *cf.Client,
	filterRegex *regexp.Regexp,
	configSource types.ConfigSource,
	tags []string,
) ([]types.Stack, error) {
	var stacks []types.Stack

	pageIndex := 0
	var nextToken *string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		input := cf.ListStacksInput{
			NextToken: nextToken,
			StackStatusFilter: []cftypes.StackStatus{
				cftypes.StackStatusCreateInProgress,
				cftypes.StackStatusCreateComplete,
				cftypes.StackStatusRollbackInProgress,
				cftypes.StackStatusRollbackFailed,
				cftypes.StackStatusRollbackComplete,
				cftypes.StackStatusDeleteFailed,
				cftypes.StackStatusUpdateInProgress,
				cftypes.StackStatusUpdateCompleteCleanupInProgress,
				cftypes.StackStatusUpdateComplete,
				cftypes.StackStatusUpdateFailed,
				cftypes.StackStatusUpdateRollbackInProgress,
				cftypes.StackStatusUpdateRollbackFailed,
				cftypes.StackStatusUpdateRollbackCompleteCleanupInProgress,
				cftypes.StackStatusUpdateRollbackComplete,
			},
		}
		output, err := client.ListStacks(ctx, &input)
		if err != nil {
			return stacks, err
		}
		for _, st := range output.StackSummaries {
			if st.StackName == nil {
				continue
			}

			if st.StackId == nil {
				continue
			}

			if filterRegex != nil && !filterRegex.Match([]byte(*st.StackName)) {
				continue
			}

			stack := types.Stack{
				Name:         *st.StackName,
				Arn:          *st.StackId,
				ConfigSource: configSource,
				TemplatePath: nil,
				Tags:         tags,
			}

			stacks = append(stacks, stack)
		}

		if output.NextToken == nil {
			break
		}

		if pageIndex >= getStacksMaxPages {
			break
		}

		nextToken = output.NextToken
	}

	return stacks, nil
}
