package releases

import (
	"blatta/pkg/db"
	"blatta/pkg/releases"
	"fmt"
	"github.com/spf13/cobra"
)

var SaveUrlFlag string

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Releases save command",
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := db.NewPoolFromUrl(SaveUrlFlag)
		if err != nil {
			panic(err)
		}
		errors := releases.UpdateReleases(pool)
		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Println(err)
			}
			panic(errors)
		}
	},
}

func init() {
	releasesCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringVar(&SaveUrlFlag, "url", "", "DB connection URL")
	saveCmd.MarkFlagRequired("url")
}
