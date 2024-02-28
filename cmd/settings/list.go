package settings

import (
	settings2 "blatta/pkg/settings"
	"fmt"
	"github.com/spf13/cobra"
)

var listSettingsVersionFlag string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Settings list command",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := settings2.ClusterSettingsFromRelease(listSettingsVersionFlag)
		if err != nil {
			panic(err)
		}
		for _, s := range settings {
			fmt.Println(s)
		}
	},
}

func init() {
	settingsCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listSettingsVersionFlag, "version", "v23.2.1", "CRDB version, starting with 'v'")
}
