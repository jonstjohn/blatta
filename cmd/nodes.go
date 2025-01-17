/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"blatta/api"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/tabwriter"
	"time"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		apiUrl := viper.GetString("url")
		username := viper.GetString("username")
		password := viper.GetString("password")
		//cacert := viper.GetString("cacert")
		//pgUrl := viper.GetString("pgurl")
		insecure := viper.GetBool("insecure")

		apiKey, err := api.Login(apiUrl, username, password, insecure)
		if err != nil {
			return err
		}

		nodes := api.GetNodes(apiUrl, apiKey, insecure)
		printNodes(nodes)

		return nil
	},
}

func printNodes(nodes []api.Node) {

	t := time.Now()
	fmt.Println(t.Format(time.RFC3339))

	w := tabwriter.NewWriter(os.Stdout, 7, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w,"Node\tCPU\t")

	for _, n := range nodes {
		fmt.Fprintf(w,"%d\t%.1f%%\t\n", n.NodeId, n.Metrics.SysCpuUserPercent*100)
	}
	w.Flush()
}

func init() {
	monitorCmd.AddCommand(nodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
