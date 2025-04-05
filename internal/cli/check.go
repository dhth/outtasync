package cli

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dhth/outtasync/internal/aws"
	"github.com/dhth/outtasync/internal/types"
)

type result struct {
	stack       types.Stack
	syncResult  *types.StackTemplateCompared
	driftResult *types.StackDriftCheckResult
}

func ShowCheckResults(
	stacks []types.Stack,
	cfClients map[string]aws.CFClient,
	checkTemplate, checkDrift, showProgressIndicator bool,
	format types.CheckOutputFormat,
	listNegativesOnly bool,
	htmlOutputConfig *types.CheckHTMLOutputConfig,
) {
	done := make(chan bool)
	templateChan := make(chan struct{})
	driftChan := make(chan struct{})

	totalCompareChecks := 0
	var stacksToCheck []types.Stack

	for _, stack := range stacks {
		cfgKey := stack.AWSConfigKey()
		_, ok := cfClients[cfgKey]
		if !ok {
			continue
		}

		if (checkTemplate && stack.TemplatePath != nil) || checkDrift {
			stacksToCheck = append(stacksToCheck, stack)
		}
		if checkTemplate && stack.TemplatePath != nil {
			totalCompareChecks++
		}
	}
	if len(stacksToCheck) == 0 {
		return
	}

	showProgressIndicator = showProgressIndicator && len(stacksToCheck) >= 1

	if showProgressIndicator {
		go showLoadingSpinner(done, templateChan, driftChan, checkTemplate, checkDrift, totalCompareChecks, len(stacksToCheck))
	}

	syncSemaphore := make(chan struct{}, 10)
	driftSemaphore := make(chan struct{}, 3)
	syncResultChannel := make(chan types.StackTemplateCompared)
	var syncWg sync.WaitGroup
	driftResultChan := make(chan types.StackDriftCheckResult)
	var driftWg sync.WaitGroup

	results := make(map[string]result)

	for _, stack := range stacksToCheck {
		results[stack.Key()] = result{stack: stack}
	}

	if checkTemplate {
		for _, stack := range stacksToCheck {
			cfgKey := stack.AWSConfigKey()
			client, ok := cfClients[cfgKey]
			if !ok {
				continue
			}

			if stack.TemplatePath == nil {
				continue
			}

			syncWg.Add(1)
			go func(stack types.Stack, cfClient aws.CFClient) {
				defer syncWg.Done()
				syncSemaphore <- struct{}{}
				defer func() {
					<-syncSemaphore
				}()
				syncResultChannel <- aws.CompareStackTemplateCode(cfClient, stack.Name, stack.Key(), *stack.TemplatePath, stack.TemplateRemoteCallHeaders)
				if showProgressIndicator {
					templateChan <- struct{}{}
				}
			}(stack, client)
		}

		go func() {
			syncWg.Wait()
			close(syncResultChannel)
		}()

		for r := range syncResultChannel {
			stackResult, ok := results[r.StackKey]
			if !ok {
				continue
			}
			stackResult.syncResult = &r
			results[r.StackKey] = stackResult
		}
	}

	if checkDrift {
		for _, stack := range stacksToCheck {
			cfgKey := stack.AWSConfigKey()
			client, ok := cfClients[cfgKey]
			if !ok {
				continue
			}

			driftWg.Add(1)
			go func(stack types.Stack, cfClient aws.CFClient) {
				defer driftWg.Done()
				driftSemaphore <- struct{}{}
				defer func() {
					<-driftSemaphore
				}()
				driftResultChan <- aws.CheckStackDriftStatus(cfClient, stack)
				if showProgressIndicator {
					driftChan <- struct{}{}
				}
			}(stack, client)
		}

		go func() {
			driftWg.Wait()
			close(driftResultChan)
		}()

		for r := range driftResultChan {
			stackResult, ok := results[r.Stack.Key()]
			if !ok {
				continue
			}
			stackResult.driftResult = &r
			results[r.Stack.Key()] = stackResult
		}
	}

	if showProgressIndicator {
		done <- true
	}
	close(templateChan)
	close(driftChan)

	if len(results) == 0 {
		return
	}

	var report string
	switch format {
	case types.Default:
		report = getDefaultReport(results, stacksToCheck, checkTemplate, checkDrift, listNegativesOnly)
	case types.Delimited:
		report = getDelimitedReport(results, stacksToCheck, checkTemplate, checkDrift, listNegativesOnly)
	case types.HTML:
		report = getHTMLReport(results, stacksToCheck, checkTemplate, checkDrift, listNegativesOnly, htmlOutputConfig)
	}

	fmt.Println(report)
}

func showLoadingSpinner(done chan bool,
	templateChecked, driftChecked chan struct{},
	showComparisonStatus, showDriftStatus bool,
	totalTemplateChecks, totalDriftChecks int,
) {
	var numTemplatesChecked int
	var numDriftChecked int
	spinner := []rune{'⠷', '⠯', '⠟', '⠻', '⠽', '⠾'}
	for {
		select {
		case <-done:
			fmt.Fprint(os.Stderr, "\r\033[K")
			return
		case <-templateChecked:
			numTemplatesChecked++
		case <-driftChecked:
			numDriftChecked++
		default:
			for _, r := range spinner {
				if showComparisonStatus && showDriftStatus {
					if totalTemplateChecks > 0 {
						fmt.Fprintf(os.Stderr, "\r\033[K%c templates checked: %d/%d; drift checked: %d/%d", r, numTemplatesChecked, totalTemplateChecks, numDriftChecked, totalDriftChecks)
					} else {
						fmt.Fprintf(os.Stderr, "\r\033[K%c no templates configured; drift checked: %d/%d", r, numDriftChecked, totalDriftChecks)
					}
				} else if showComparisonStatus {
					fmt.Fprintf(os.Stderr, "\r\033[K%c templates checked: %d/%d", r, numTemplatesChecked, totalTemplateChecks)
				} else {
					fmt.Fprintf(os.Stderr, "\r\033[K%c drift checked: %d/%d", r, numDriftChecked, totalDriftChecks)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
