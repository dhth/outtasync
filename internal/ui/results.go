package ui

import (
	"fmt"
	"os"
)

func ShowResults(stacks []Stack, awsCfgs map[string]AwsConfig) {
	results := make(map[string]StackSyncResult)
	resultChannel := make(chan StackSyncResult)

	for _, stack := range stacks {
		cfgKey := GetAWSConfigKey(stack)
		go func(stack Stack, awsCfg AwsConfig) {
			resultChannel <- CheckStackSyncStatus(awsCfg, stack)
		}(stack, awsCfgs[cfgKey])
	}

	for range stacks {
		r := <-resultChannel
		results[r.Stack.AwsProfile+":"+r.Stack.AwsRegion+":"+r.Stack.Name] = r
	}

	var outtaSyncStacks string
	var outtaSyncCount int

	var errorResults string
	var errorsCount int
	for k, r := range results {
		if r.Err != nil {
			errorResults += fmt.Sprintf("%s %s\n", RightPadTrim(k, 80), r.Err.Error())
			errorsCount++
		} else {
			if r.Outtasync {
				outtaSyncStacks += k + "\n"
				outtaSyncCount++
			}
		}
	}

	if outtaSyncCount > 0 {
		if outtaSyncCount == 1 {
			fmt.Printf("1 stack is outtasync:\n\n")
		} else {
			fmt.Printf(fmt.Sprintf("%d stacks are outtasync:\n\n", outtaSyncCount))
		}
		fmt.Print(outtaSyncStacks)
	}

	if errorsCount > 0 {
		fmt.Println("")
		if errorsCount == 1 {
			fmt.Printf("1 error:\n\n")
		} else {
			fmt.Printf(fmt.Sprintf("%d errors:\n\n", errorsCount))
		}
		fmt.Print(errorResults)
	}

	if outtaSyncCount > 0 || errorsCount > 0 {
		os.Exit(1)
	}
}
