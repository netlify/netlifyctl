package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/operations"
	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:   "netlify",
		Short: "Command Line Interface for netlify.com",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			logrus.WithFields(logrus.Fields{"command": cmd.Use, "arguments": args}).Debug("PreRun")
			if debug {
				cmd.DebugFlags()
			}
		},
	}
)

// Execute configures all the commands and runs the root.
func Execute() {
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().StringP("endpoint", "E", "https://api.netlify.com", "default API endpoint")
	rootCmd.PersistentFlags().StringP("streaming", "S", "wss://streaming.netlify.com", "default streaming API endpoint")
	rootCmd.PersistentFlags().StringP("access_token", "A", "", "access token for Netlify's API")

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "enable debug tracing")
	rootCmd.PersistentFlags().BoolVarP(&operations.AssumeYes, "yes", "y", false, "automatic yes to confirmation prompts")

	err := configuration.SetupViper(rootCmd.PersistentFlags(), rootCmd.Flags())
	if err != nil {
		fmt.Println("error while configuring CLI: " + err.Error())
		os.Exit(1)
	}

	addCommands()

	if c, err := rootCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println(c.UsageString())
		}

		os.Exit(-1)
	}
}

var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}
