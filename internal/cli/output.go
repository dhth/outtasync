package cli

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/dhth/outtasync/internal/types"
)

//go:embed assets/template.html
var defaultHTMLTemplate string

const (
	err = "error"
	no  = "no"
	yes = "yes"
)

func getDefaultReport(
	results map[string]result,
	stacks []types.Stack,
	showTemplateResults,
	showDriftResults bool,
	listNegativesOnly bool,
) string {
	var errors []string
	rows := make([][]string, len(stacks))
	for i, stack := range stacks {
		result, ok := results[stack.Key()]
		var row []string
		row = append(row, stack.Name)

		if !ok {
			errors = append(errors, "unexpected error")
			if showTemplateResults {
				if result.syncResult == nil {
					row = append(row, errorStyle.Render(""))
				} else {
					row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
				}
			}
			if showDriftResults {
				if result.driftResult == nil {
					row = append(row, errorStyle.Render(""))
				} else {
					row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
				}
			}
			rows[i] = row
			continue
		}

		isNegative := false

		if showTemplateResults {
			if result.syncResult == nil {
				row = append(row, errorStyle.Render(""))
			} else {
				if result.syncResult.Err != nil {
					errors = append(errors, result.syncResult.Err.Error())
					row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
				} else {
					if result.syncResult.Mismatch {
						row = append(row, outtaSyncStyle.Render(no))
						isNegative = true
					} else {
						row = append(row, inSyncStyle.Render(yes))
					}
				}
			}
		}

		if showDriftResults {
			if result.driftResult == nil {
				row = append(row, errorStyle.Render(""))
			} else {
				if result.driftResult.Err != nil {
					errors = append(errors, result.driftResult.Err.Error())
					row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
				} else {
					if result.driftResult.Output.DetectionStatus == cftypes.StackDriftDetectionStatusDetectionComplete {
						if result.driftResult.Output.StackDriftStatus == cftypes.StackDriftStatusDrifted {
							row = append(row, outtaSyncStyle.Render(no))
							isNegative = true
						} else {
							row = append(row, inSyncStyle.Render(yes))
						}
					} else {
						errors = append(errors, fmt.Sprintf("drift detection could not be completed, detection status: %v", result.driftResult.Output.DetectionStatus))
						row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
					}
				}
			}
		}

		if listNegativesOnly && !isNegative {
			continue
		}

		rows[i] = row
	}

	if len(rows) == 0 {
		return ""
	}

	//nolint:prealloc
	var outputLines []string
	if showTemplateResults && showDriftResults {
		outputLines = append(outputLines,
			fmt.Sprintf("%s\t%s\t%s",
				stackHeadingStyle.Render("stack"),
				statusHeadingStyle.Render("template-in-sync"),
				statusHeadingStyle.Render("no-drift"),
			),
		)
	} else if showTemplateResults {
		outputLines = append(outputLines,
			fmt.Sprintf("%s\t%s",
				stackHeadingStyle.Render("stack"),
				statusHeadingStyle.Render("template-in-sync"),
			),
		)
	} else {
		outputLines = append(outputLines,
			fmt.Sprintf("%s\t%s",
				stackHeadingStyle.Render("stack"),
				statusHeadingStyle.Render("no-drift"),
			),
		)
	}

	for _, r := range rows {
		if len(r) == 3 {
			outputLines = append(outputLines, fmt.Sprintf("%s\t%s\t%s", stackNameStyle.Render(r[0]), r[1], r[2]))
		} else if len(r) == 2 {
			outputLines = append(outputLines, fmt.Sprintf("%s\t%s", stackNameStyle.Render(r[0]), r[1]))
		}
	}

	if len(errors) > 0 {
		outputLines = append(outputLines, "")
		outputLines = append(outputLines, errorStyle.Render("Errors"))
		for i, e := range errors {
			outputLines = append(outputLines, fmt.Sprintf("[%d]: %s", i+1, e))
		}
	}

	return strings.Join(outputLines, "\n")
}

func getDelimitedReport(
	results map[string]result,
	stacks []types.Stack,
	showTemplateResults,
	showDriftResults bool,
	listNegativesOnly bool,
) string {
	rows := make([][]string, len(stacks))
	for i, stack := range stacks {
		result, ok := results[stack.Key()]
		var row []string
		row = append(row, stack.Name)
		if !ok {
			if showTemplateResults {
				row = append(row, err)
			}
			if showDriftResults {
				row = append(row, err)
			}
			rows[i] = row
			continue
		}

		isNegative := false

		if showTemplateResults {
			if result.syncResult == nil {
				row = append(row, "")
			} else {
				if result.syncResult.Err != nil {
					row = append(row, err)
					isNegative = true
				} else {
					if result.syncResult.Mismatch {
						row = append(row, no)
						isNegative = true
					} else {
						row = append(row, yes)
					}
				}
			}
		}

		if showDriftResults {
			if result.driftResult == nil {
				row = append(row, "")
			} else {
				if result.driftResult.Err != nil {
					row = append(row, err)
					isNegative = true
				} else {
					if result.driftResult.Output.DetectionStatus == cftypes.StackDriftDetectionStatusDetectionComplete {
						if result.driftResult.Output.StackDriftStatus == cftypes.StackDriftStatusDrifted {
							row = append(row, no)
							isNegative = true
						} else {
							row = append(row, yes)
						}
					} else {
						row = append(row, err)
					}
				}
			}
		}

		if listNegativesOnly && !isNegative {
			continue
		}

		rows[i] = row
	}

	if len(rows) == 0 {
		return ""
	}

	//nolint:prealloc
	var outputLines []string
	if showTemplateResults && showDriftResults {
		outputLines = append(outputLines, "stack,template-in-sync,no-drift")
	} else if showTemplateResults {
		outputLines = append(outputLines, "stack,template-in-sync")
	} else {
		outputLines = append(outputLines, "stack,no-drift")
	}
	for _, r := range rows {
		outputLines = append(outputLines, strings.Join(r, ","))
	}

	return strings.Join(outputLines, "\n")
}

func getHTMLReport(
	results map[string]result,
	stacks []types.Stack,
	showTemplateResults,
	showDriftResults bool,
	listNegativesOnly bool,
	htmlOutputConfig *types.CheckHTMLOutputConfig,
) string {
	var errors []string
	rows := make([]HTMLDataRow, len(stacks))
	for i, stack := range stacks {
		result, ok := results[stack.Key()]
		row := HTMLDataRow{StackName: stack.Name}

		if !ok {
			errors = append(errors, fmt.Sprintf("[%d] unexpected error", len(errors)+1))
			if showTemplateResults {
				if result.syncResult == nil {
					row.TemplateCheckValue = "-"
				} else {
					row.TemplateCheckValue = fmt.Sprintf("error [%d]", len(errors))
					row.TemplateCheckErrored = true
				}
			}
			if showDriftResults {
				if result.driftResult == nil {
					row.DriftCheckValue = "-"
				} else {
					row.DriftCheckValue = fmt.Sprintf("error [%d]", len(errors))
					row.DriftCheckError = true
				}
			}
			rows[i] = row
			continue
		}

		isNegative := false

		if showTemplateResults {
			if result.syncResult == nil {
				row.TemplateCheckValue = "-"
			} else {
				if result.syncResult.Err != nil {
					errors = append(errors, fmt.Sprintf("[%d] %s", len(errors)+1, result.syncResult.Err.Error()))
					row.TemplateCheckValue = fmt.Sprintf("error [%d]", len(errors))
					row.TemplateCheckErrored = true
				} else {
					if result.syncResult.Mismatch {
						row.TemplateCheckValue = no
						isNegative = true
					} else {
						row.TemplateCheckValue = yes
						row.TemplateInSync = true
					}
				}
			}
		}

		if showDriftResults {
			if result.driftResult == nil {
				row.DriftCheckValue = ""
			} else {
				if result.driftResult.Err != nil {
					errors = append(errors, fmt.Sprintf("[%d] %s", len(errors)+1, result.driftResult.Err.Error()))
					row.DriftCheckValue = fmt.Sprintf("error [%d]", len(errors))
					row.DriftCheckError = true
				} else {
					if result.driftResult.Output.DetectionStatus == cftypes.StackDriftDetectionStatusDetectionComplete {
						if result.driftResult.Output.StackDriftStatus == cftypes.StackDriftStatusDrifted {
							row.DriftCheckValue = no
							row.NoDrift = false
							isNegative = true
						} else {
							row.DriftCheckValue = yes
							row.NoDrift = true
						}
					} else {
						errors = append(errors, fmt.Sprintf("[%d] drift detection could not be completed, detection status: %v", len(errors)+1, result.driftResult.Output.DetectionStatus))
						row.DriftCheckValue = fmt.Sprintf("error [%d]", len(errors))
						row.DriftCheckError = true
					}
				}
			}
		}

		if listNegativesOnly && !isNegative {
			continue
		}

		rows[i] = row
	}

	if len(rows) == 0 {
		return ""
	}

	var title string
	var htmlTemplate string
	if htmlOutputConfig != nil {
		title = htmlOutputConfig.Title
		if htmlOutputConfig.Template != "" {
			htmlTemplate = htmlOutputConfig.Template
		} else {
			htmlTemplate = defaultHTMLTemplate
		}
	}

	htmlData := CheckReportHTMLData{
		Title:               title,
		Timestamp:           time.Now().Format("2006-01-02 15:04:05 MST"),
		Rows:                rows,
		Errors:              errors,
		ShowTemplateResults: showTemplateResults,
		ShowDriftResults:    showDriftResults,
	}

	var tmpl *template.Template
	var err error
	tmpl, err = template.New("outtasync").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, htmlData)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}

	return buf.String()
}

type CheckReportHTMLData struct {
	Title               string
	Timestamp           string
	Rows                []HTMLDataRow
	Errors              []string
	ShowTemplateResults bool
	ShowDriftResults    bool
}

type HTMLDataRow struct {
	StackName            string
	TemplateCheckValue   string
	TemplateInSync       bool
	TemplateCheckErrored bool
	DriftCheckValue      string
	NoDrift              bool
	DriftCheckError      bool
}
