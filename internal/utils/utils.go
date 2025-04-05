package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func ExpandTilde(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func GetDiff(template, actual string) ([]byte, error) {
	var zero []byte
	tmpFileTemplate, err := os.CreateTemp("", "outtasync-*.yml")
	defer func() {
		_ = tmpFileTemplate.Close()
	}()
	if err != nil {
		return zero, err
	}

	_, err = tmpFileTemplate.WriteString(template)
	if err != nil {
		return zero, err
	}

	tmpFileActual, err := os.CreateTemp("", "outtasync-*.yml")
	defer func() {
		_ = tmpFileActual.Close()
	}()
	if err != nil {
		return zero, err
	}

	_, err = tmpFileActual.WriteString(actual)
	if err != nil {
		return zero, err
	}

	c := exec.Command("bash", "-c",
		fmt.Sprintf("git --no-pager diff --src-prefix='Template ' --dst-prefix='Actual ' --no-index -- %s %s || true",
			tmpFileTemplate.Name(),
			tmpFileActual.Name(),
		))
	out, err := c.Output()
	if err != nil {
		return zero, err
	}

	return out, nil
}
