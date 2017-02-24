package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	aoperations "github.com/netlify/open-api/go/plumbing/operations"

	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/operations"
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
	rootCmd.SilenceErrors = true

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "E", "https://api.netlify.com", "default API endpoint")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", globalConfigFileName, "configuration file")
	rootCmd.PersistentFlags().StringVarP(&auth.AccessToken, "access-token", "A", "", "access token for Netlify's API")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "enable debug tracing")

	rootCmd.PersistentFlags().BoolVarP(&operations.AssumeYes, "yes", "y", false, "automatic yes to confirmation prompts")

	addCommands()
	if c, err := rootCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println(c.UsageString())
		} else if uErr, ok := err.(*aoperations.ListSitesDefault); ok {
			errStr := fmt.Sprintf("%d", uErr.Code())
			if uErr.Payload.Message != nil {
				errStr += ": " + *uErr.Payload.Message
			}
			fmt.Fprintf(os.Stderr, "%s\n", errStr)
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
