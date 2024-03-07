package releases

import (
	"blatta/pkg/releases"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Releases list command",
	Run: func(cmd *cobra.Command, args []string) {
		rp := releases.NewRemoteDataSource()
		releases, err := rp.GetReleases()
		if err != nil {
			panic(err)
		}
		b, err := json.MarshalIndent(releases, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	},
}

func init() {
	releasesCmd.AddCommand(listCmd)
}
