package types

import "fmt"

type Stack struct {
	Name           string
	AwsProfile     string
	AwsRegion      string
	Local          string
	Template       string
	Tags           []string
	RefreshCommand string
}

type StackSyncResult struct {
	Stack        Stack
	TemplateBody string
	Outtasync    bool
	Err          error
}

func (stack Stack) Key() string {
	return fmt.Sprintf("%s:%s:%s", stack.AwsProfile, stack.AwsRegion, stack.Name)
}

func (stack Stack) AWSConfigKey() string {
	return stack.AwsProfile + ":" + stack.AwsRegion
}
