package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const options = " (yes/no) "

var AssumeYes bool

func Error(err error) {
	fmt.Errorf("%v", err)
}

func Warning(warn string) {
	fmt.Printf("[WARNING] %v\n", warn)
}

func AskForConfirmation(message string) bool {
	if AssumeYes {
		return true
	}

	var response string

	fmt.Print(message + options)
	_, err := fmt.Scanln(&response)

	if err != nil {
		return false
	}

	switch {
	case response[0] == 'y' || response[0] == 'Y':
		return true
	case response[0] == 'n' || response[0] == 'N':
		return false
	default:
		return AskForConfirmation(message)
	}
}

func AskForInput(message, defaultValue string) (string, error) {
	if AssumeYes && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	if len(defaultValue) > 0 {
		fmt.Printf("%s (default: %s) ", message, defaultValue)
	} else {
		fmt.Printf("%s ", message)
	}

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return response, err
	}

	response = strings.TrimSpace(response)

	if len(response) == 0 && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	if response == defaultValue {
		return response, nil
	}

	return response, nil
}

func DoneCheck() string {
	return color.GreenString("✔")
}

func ErrorCheck() string {
	return color.RedString("✘")
}

func WorldCheck() string {
	return color.GreenString("🌎")
}

func Bold(msg string, args ...interface{}) {
	c := color.New(color.Bold)
	c.Printf(msg, args...)
}

func ConfirmCreateSite() bool {
	return AskForConfirmation("We cannot find your site, do you want to create a new one?")
}

func ConfirmOverwriteSite() bool {
	return AskForConfirmation("There's already a site ID stored for this folder. Ignore and create a new site?")
}

func Track(process, success string, task func() error) error {
	tt := NewTaskTracker()
	return TrackWithTracker(process, success, tt, task)
}

func TrackWithTracker(process, success string, tt *TaskTracker, task func() error) (err error) {
	defer func() {
		if err != nil {
			tt.Failure(process)
		} else {
			tt.Success(success)
		}
	}()

	tt.Start(process)
	err = task()
	return
}
