package settings

import (
	settings2 "blatta/pkg/settings"
	"github.com/spf13/cobra"
)

var saveSettingsVersionFlag string
var saveSettingsUrlFlag string

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Settings save command",
	Run: func(cmd *cobra.Command, args []string) {
		err := settings2.SaveClusterSettingsForVersion(saveSettingsVersionFlag, saveSettingsUrlFlag)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	settingsCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringVar(&saveSettingsVersionFlag, "version", "v23.2.1", "Specify a single CRDB version, starting with 'v'")
	saveCmd.Flags().StringVar(&saveSettingsUrlFlag, "url", "", "DB connection URL")
	saveCmd.MarkFlagRequired("url")
}
