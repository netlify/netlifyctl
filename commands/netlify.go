package commands

import (
	"os"
	"regexp"

	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/operations"
	"github.com/spf13/cobra"
)

const globalConfigFileName = "netlify.toml"

var (
	configFile string
	debug      bool
	endpoint   string

	rootCmd = &cobra.Command{
		Use:   "netlify",
		Short: "Command Line Interface for netlify.com",
	}
)

// Execute configures all the commands and runs the root.
func Execute() {
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "E", "https://api.netlify.com", "default API endpoint")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", globalConfigFileName, "configuration file")
	rootCmd.PersistentFlags().StringVarP(&auth.AccessToken, "access-token", "A", "", "access token for Netlify's API")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "enable debug tracing")

	rootCmd.PersistentFlags().BoolVarP(&operations.AssumeYes, "yes", "y", false, "automatic yes to confirmation prompts")

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
