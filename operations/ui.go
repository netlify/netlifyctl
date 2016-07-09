package operations

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const options = " (yes/no) "

func askForConfirmation(message string) bool {
	var response string
	fmt.Print("=> " + message + options)
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
		fmt.Println("=> Please type `yes` or `no` and then press enter")
		return askForConfirmation(message)
	}
}

func askForInput(message, defaultValue string, validate validator) (string, error) {
	if len(defaultValue) > 0 {
		fmt.Printf("=> %s (default: %s) ", message, defaultValue)
	} else {
		fmt.Printf("=> %s ", message)
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

	if err := validate(response); err != nil {
		fmt.Println(err)
		return askForInput(message, defaultValue, validate)
	}

	return response, nil
}
