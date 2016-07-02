package commands

import (
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/netlify/netlifyctl/auth"
	"github.com/spf13/cobra"
)

var (
	debug   bool
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
	rootCmd.PersistentFlags().StringVarP(&auth.AccessToken, "access-token", "A", "", "access token for Netlify's API")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "enable debug tracing")

	addCommands()

	if c, err := rootCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println(c.UsageString())
		}

		os.Exit(-1)
	}
}

func addCommands() {
	rootCmd.AddCommand(sitesCmd)
	rootCmd.AddCommand(deployCmd)
}

var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}
