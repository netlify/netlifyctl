package errors

import "fmt"

type CommandError struct {
	s         string
	userError bool
}

func (u CommandError) Error() string {
	return u.s
}

func (u CommandError) IsUserError() bool {
	return u.userError
}

// ArgumentError is for when there is a problem with the args passed in the CLI
type ArgumentError struct {
	errMsg string
}

func (e ArgumentError) Error() string {
	return e.errMsg
}

func UserError(a ...interface{}) CommandError {
	return CommandError{s: fmt.Sprintln(a...), userError: true}
}

func UserErrorF(format string, a ...interface{}) CommandError {
	return CommandError{s: fmt.Sprintf(format, a...), userError: true}
}

func SystemError(a ...interface{}) CommandError {
	return CommandError{s: fmt.Sprintln(a...), userError: false}
}

func SystemErrorF(format string, a ...interface{}) CommandError {
	return CommandError{s: fmt.Sprintf(format, a...), userError: false}
}

func ArgumentErrorF(format string, a ...interface{}) ArgumentError {
	return ArgumentError{errMsg: fmt.Sprintf(format, a...)}
}
