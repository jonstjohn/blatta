package releases

import (
	"blatta/pkg/dbpgx"
	"blatta/pkg/releases"
	"github.com/spf13/cobra"
)

var SaveUrlFlag string

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Releases save command",
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := dbpgx.NewPoolFromUrl(SaveUrlFlag)
		if err != nil {
			panic(err)
		}
		err = releases.UpdateReleases(pool)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	releasesCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringVar(&SaveUrlFlag, "url", "", "DB connection URL")
	saveCmd.MarkFlagRequired("url")
}
