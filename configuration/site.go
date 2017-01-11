package configuration

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SiteIdForCommand(cmd *cobra.Command) (string, error) {
	var siteId string
	if siteIdFlag := cmd.Flag("site-id"); siteIdFlag != nil {
		siteId = siteIdFlag.Value.String()
	}

	var configFile = cmd.Root().Flag("config").Value.String()
	if siteId == "" && Exist(configFile) {
		conf, err := Load(configFile)
		if err != nil {
			return "", err
		}

		siteId = conf.Settings.ID
	}

	if siteId == "" {
		return "", fmt.Errorf("missing site ID")
	}

	return siteId, nil
}
