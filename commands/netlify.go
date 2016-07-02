package commands

import (
	"os"
	"regexp"

	"github.com/netlify/netlify-go-cli/auth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netlify",
	Short: "CLI interface to netlify.com",
}

// Execute configures all the commands and runs the root.
func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().StringVarP(&auth.AccessToken, "access-token", "A", "", "access token for Netlify's API")

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
}

var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}
