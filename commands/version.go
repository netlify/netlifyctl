package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const unknown = "unknown"

var (
	SHA     string
	Version string
)

var versionCmd = &cobra.Command{
	Run:     showVersion,
	Use:     "version",
	Aliases: []string{"v"},
}

func showVersion(cmd *cobra.Command, args []string) {
	v := Version
	if v == "" {
		v = unknown
	}
	s := SHA
	if s == "" {
		s = unknown
	}
	fmt.Printf("Version: %s\nGit SHA: %s\n", v, s)
}
