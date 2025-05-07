package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dhth/outtasync/internal/constants"
	"github.com/dhth/outtasync/internal/types"

	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

const (
	driftCheckTimeoutSecs = 15
	driftCheckMaxAttempts = 5
)

var (
	errDriftDetectionTimedOut = errors.New("drift detection timed out")
	errDriftDetectionFailed   = errors.New("drift detection failed")
)

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

		switch descOutput.DetectionStatus {
		case cftypes.StackDriftDetectionStatusDetectionComplete:
			return types.StackDriftCheckResult{
				Stack:  stack,
				Output: descOutput,
			}
		case cftypes.StackDriftDetectionStatusDetectionFailed:
			if descOutput.DetectionStatusReason != nil {
				return types.StackDriftCheckResult{
					Stack: stack,
					Err:   fmt.Errorf("%w; reason: %s", errDriftDetectionFailed, *descOutput.DetectionStatusReason),
				}
			}
			return types.StackDriftCheckResult{
				Stack: stack,
				Err:   errDriftDetectionFailed,
			}
		default:
			time.Sleep(time.Millisecond * describeDriftSleepMillis)
		}
	}

	return types.StackDriftCheckResult{
		Stack: stack,
		Err: fmt.Errorf("%w; if you think the deadline for this check should be higher, let %s know via %s",
			errDriftDetectionTimedOut,
			constants.Author,
			constants.RepoIssuesURL),
	}
}
