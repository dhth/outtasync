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
	//nolint:prealloc
	var rows [][]string
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
			row[i] = stackErrorResultStyle.Render(row[0])
			rows = append(rows, row)

			continue
		}

		hasError := false
		hasNegativeResult := false

		if showTemplateResults {
			if result.syncResult == nil {
				row = append(row, errorStyle.Render(""))
			} else {
				if result.syncResult.Err != nil {
					errors = append(errors, result.syncResult.Err.Error())
					row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
					hasError = true
				} else {
					if result.syncResult.Mismatch {
						row = append(row, outtaSyncStyle.Render(no))
						hasNegativeResult = true
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
					hasError = true
				} else {
					if result.driftResult.Output.DetectionStatus == cftypes.StackDriftDetectionStatusDetectionComplete {
						if result.driftResult.Output.StackDriftStatus == cftypes.StackDriftStatusDrifted {
							row = append(row, outtaSyncStyle.Render(no))
							hasNegativeResult = true
						} else {
							row = append(row, inSyncStyle.Render(yes))
						}
					} else {
						errors = append(errors, fmt.Sprintf("drift detection could not be completed, detection status: %v", result.driftResult.Output.DetectionStatus))
						row = append(row, errorStyle.Render(fmt.Sprintf("error [%d]", len(errors))))
						hasError = true
					}
				}
			}
		}

		if listNegativesOnly && !hasNegativeResult {
			continue
		}

		if hasNegativeResult {
			row[0] = stackNegativeResultStyle.Render(row[0])
		} else if hasError {
			row[0] = stackErrorResultStyle.Render(row[0])
		} else {
			row[0] = stackPositiveResultStyle.Render(row[0])
		}

		rows = append(rows, row)
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
				statusStyle.Render("template-in-sync"),
				statusStyle.Render("no-drift"),
			),
		)
	} else if showTemplateResults {
		outputLines = append(outputLines,
			fmt.Sprintf("%s\t%s",
				stackHeadingStyle.Render("stack"),
				statusStyle.Render("template-in-sync"),
			),
		)
	} else {
		outputLines = append(outputLines,
			fmt.Sprintf("%s\t%s",
				stackHeadingStyle.Render("stack"),
				statusStyle.Render("no-drift"),
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
	//nolint:prealloc
	var rows [][]string
	for _, stack := range stacks {
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
			rows = append(rows, row)
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

		rows = append(rows, row)
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
) (string, error) {
	var errors []string
	var diffs []HTMLDiff
	//nolint:prealloc
	var rows []HTMLDataRow
	for _, stack := range stacks {
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
					row.HasError = true
				}
			}
			if showDriftResults {
				if result.driftResult == nil {
					row.DriftCheckValue = "-"
				} else {
					row.DriftCheckValue = fmt.Sprintf("error [%d]", len(errors))
					row.DriftCheckError = true
					row.HasError = true
				}
			}
			rows = append(rows, row)
			continue
		}

		if showTemplateResults {
			if result.syncResult == nil {
				row.TemplateCheckValue = "-"
			} else {
				if result.syncResult.Err != nil {
					errors = append(errors, fmt.Sprintf("[%d] %s", len(errors)+1, result.syncResult.Err.Error()))
					row.TemplateCheckValue = fmt.Sprintf("error [%d]", len(errors))
					row.TemplateCheckErrored = true
					row.HasError = true
				} else {
					if result.syncResult.Mismatch {
						row.TemplateCheckValue = no
						row.HasNegativeResult = true
						if result.syncResult.DiffErr != nil {
							diffs = append(diffs, HTMLDiff{
								StackName: stack.Name,
								Diff:      result.syncResult.DiffErr.Error(),
							})
						} else {
							diffs = append(diffs, HTMLDiff{
								StackName: stack.Name,
								Diff:      string(result.syncResult.Diff),
							})
						}
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
					row.HasError = true
				} else {
					if result.driftResult.Output.DetectionStatus == cftypes.StackDriftDetectionStatusDetectionComplete {
						if result.driftResult.Output.StackDriftStatus == cftypes.StackDriftStatusDrifted {
							row.DriftCheckValue = no
							row.NoDrift = false
							row.HasNegativeResult = true
						} else {
							row.DriftCheckValue = yes
							row.NoDrift = true
						}
					} else {
						errors = append(errors, fmt.Sprintf("[%d] drift detection could not be completed, detection status: %v", len(errors)+1, result.driftResult.Output.DetectionStatus))
						row.DriftCheckValue = fmt.Sprintf("error [%d]", len(errors))
						row.DriftCheckError = true
						row.HasError = true
					}
				}
			}
		}

		if listNegativesOnly && !row.HasNegativeResult {
			continue
		}

		rows = append(rows, row)
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
		Diffs:               diffs,
		Errors:              errors,
		ShowTemplateResults: showTemplateResults,
		ShowDriftResults:    showDriftResults,
	}

	var tmpl *template.Template
	var err error
	tmpl, err = template.New("outtasync").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, htmlData)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}

	return buf.String(), nil
}

type CheckReportHTMLData struct {
	Title               string
	Timestamp           string
	Rows                []HTMLDataRow
	Errors              []string
	Diffs               []HTMLDiff
	ShowTemplateResults bool
	ShowDriftResults    bool
}

type HTMLDataRow struct {
	StackName            string
	HasNegativeResult    bool
	HasError             bool
	TemplateCheckValue   string
	TemplateInSync       bool
	TemplateCheckErrored bool
	DriftCheckValue      string
	NoDrift              bool
	DriftCheckError      bool
}

type HTMLDiff struct {
	StackName string
	Diff      string
}
