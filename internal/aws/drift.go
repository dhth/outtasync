package aws

import (
	"context"
	"errors"
	"time"

	"github.com/dhth/outtasync/internal/types"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

const driftCheckTimeoutSecs = 15

func CheckStackDriftStatus(
	cfClient CFClient,
	stack types.Stack,
) types.StackDriftCheckResult {
	if cfClient.Err != nil {
		return types.StackDriftCheckResult{
			Stack: stack,
			Err:   cfClient.Err,
		}
	}

	client := cfClient.Client

	detectDriftParams := cf.DetectStackDriftInput{
		StackName: &stack.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), driftCheckTimeoutSecs*time.Second)
	defer cancel()

	detectOutput, err := client.DetectStackDrift(ctx, &detectDriftParams)
	if err != nil {
		return types.StackDriftCheckResult{
			Stack: stack,
			Err:   err,
		}
	}
	driftDetectionID := detectOutput.StackDriftDetectionId

	for range driftCheckMaxAttempts {
		descInput := cf.DescribeStackDriftDetectionStatusInput{
			StackDriftDetectionId: driftDetectionID,
		}
		descOutput, err := client.DescribeStackDriftDetectionStatus(ctx, &descInput)
		if err != nil {
			return types.StackDriftCheckResult{
				Stack: stack,
				Err:   err,
			}
		}

		if descOutput.DetectionStatus != cftypes.StackDriftDetectionStatusDetectionInProgress {
			return types.StackDriftCheckResult{
				Stack:  stack,
				Output: descOutput,
				Err:    err,
			}
		}
		time.Sleep(time.Millisecond * describeDriftSleepMillis)
	}

	return types.StackDriftCheckResult{
		Stack: stack,
		Err:   errors.New("couldn't fetch drift status"),
	}
}
