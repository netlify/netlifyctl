package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	aoperations "github.com/netlify/open-api/go/plumbing/operations"

	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/ui"
)

var (
	configFile string
	dump       bool
	endpoint   string

	rootCmd = &cobra.Command{
		Use:   "netlifyctl",
		Short: "Command Line Interface for netlify.com",
	}
)

// Execute configures all the commands and runs the root.
func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "E", "https://api.netlify.com", "default API endpoint")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", configuration.DefaultConfigFileName, "configuration file")
	rootCmd.PersistentFlags().StringVarP(&auth.AccessToken, "access-token", "A", "", "access token for Netlify's API")
	rootCmd.PersistentFlags().BoolVarP(&dump, "debug", "D", false, "dump debug tracing, even if there are no errors")

	rootCmd.PersistentFlags().BoolVarP(&ui.AssumeYes, "yes", "y", false, "automatic yes to confirmation prompts")

	addCommands()
	if c, err := rootCmd.ExecuteC(); err != nil {
		displayError(c, err)
		os.Exit(-1)
	}
}

var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func displayError(c *cobra.Command, raw error) {
	switch err := raw.(type) {
	case commandError:
		if err.isUserError() && userErrorRegexp.MatchString(err.Error()) {
			c.Println(c.UsageString())
		}
	case *aoperations.ListSitesDefault:
		errStr := fmt.Sprintf("%d", err.Code())
		if err.Payload.Message != "" {
			errStr += ": " + err.Payload.Message
		}
		fmt.Fprintf(os.Stderr, "%s\n", errStr)
	default:
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
