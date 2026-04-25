package app

import (
	"errors"
	"io"
	"strings"

	internalcli "github.com/HidemaruOwO/pummit/internal/cli"
	legacycli "github.com/HidemaruOwO/pummit/legacy/cli"
)

const (
	ExitCodeSuccess = 0
	ExitCodeError   = 1
)

type exitCoder interface {
	ExitCode() int
}

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func (e *ExitError) ExitCode() int {
	if e == nil {
		return ExitCodeSuccess
	}

	return e.Code
}

func NewExitError(code int, err error) error {
	if err == nil {
		return nil
	}

	return &ExitError{Code: code, Err: err}
}

func ExitCode(err error) int {
	if err == nil {
		return ExitCodeSuccess
	}

	var coded exitCoder
	if errors.As(err, &coded) {
		return coded.ExitCode()
	}

	return ExitCodeError
}

func Run(args []string, stdout, stderr io.Writer) error {
	if shouldUseNewCLI(args) {
		return internalcli.Execute(args, stdout, stderr)
	}

	return legacycli.ExecuteArgs(args)
}

func shouldUseNewCLI(args []string) bool {
	if len(args) == 0 {
		return true
	}

	if args[0] == "config" {
		return true
	}

	if args[0] == "commit" || args[0] == "alias" || args[0] == "doctor" || args[0] == "migrate" || args[0] == "mcp" {
		return true
	}

	if len(args) >= 2 {
		return true
	}

	for _, arg := range args {
		switch arg {
		case "-v", "--version", "version", "help", "-h", "--help":
			return true
		}

		if strings.HasPrefix(arg, "-") {
			continue
		}

		return false
	}

	return true
}
