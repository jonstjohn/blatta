package settings

import (
	settings2 "blatta/pkg/settings"
	"github.com/spf13/cobra"
)

var summarizeUrlFlag string

var summarizeCmd = &cobra.Command{
	Use:   "save",
	Short: "Settings save command",
	Run: func(cmd *cobra.Command, args []string) {
		err := settings2.SummarizeAndSaveSettings(summarizeUrlFlag)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	settingsCmd.AddCommand(summarizeCmd)
	summarizeCmd.Flags().StringVar(&summarizeUrlFlag, "url", "", "DB connection URL")
	summarizeCmd.MarkFlagRequired("url")
}
