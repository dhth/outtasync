package cmd

import (
	"fmt"
	"os"

	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

func showResults(stacks []types.Stack, awsCfgs map[string]aws.Config) {
	resultChannel := make(chan types.StackSyncResult)

	for _, stack := range stacks {
		cfgKey := stack.AWSConfigKey()
		go func(stack types.Stack, awsCfg aws.Config) {
			resultChannel <- aws.CheckStackSyncStatus(awsCfg, stack)
		}(stack, awsCfgs[cfgKey])
	}

	var outtaSync []string

	var errors []string

	for range stacks {
		r := <-resultChannel
		if r.Err != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", r.Stack.Key(), r.Err.Error()))
		} else if r.Outtasync {
			outtaSync = append(outtaSync, r.Stack.Key())
		}
	}

	if len(outtaSync) > 0 {
		if len(outtaSync) == 1 {
			fmt.Printf("1 stack is outtasync\n\n")
		} else {
			fmt.Printf("%d stacks are outtasync\n\n", len(outtaSync))
		}
		for _, st := range outtaSync {
			fmt.Println(st)
		}
	}

	if len(errors) > 0 {
		fmt.Println("")
		if len(errors) == 1 {
			fmt.Printf("1 error\n\n")
		} else {
			fmt.Printf("%d errors\n\n", len(errors))
		}
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	if len(outtaSync) > 0 || len(errors) > 0 {
		os.Exit(1)
	}
}
