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
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("monitor called")
	},
}

var ApiUrl string
var PgUrl string
var Username string
var Password string
var Insecure bool
var Count int
var Wait int

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")
	monitorCmd.PersistentFlags().StringVar(&ApiUrl, "url", "", "")
	monitorCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "USERNAME")
	monitorCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "PASSWORD")
	monitorCmd.PersistentFlags().BoolVar(&Insecure, "insecure", false, "Skip TLS certificate verification")
	monitorCmd.PersistentFlags().StringVar(&PgUrl, "pgurl", "", "")
	monitorCmd.PersistentFlags().IntVarP(&Count, "count", "c", 0, "Number of iterations to run (0=unlimited)")
	monitorCmd.PersistentFlags().IntVarP(&Wait, "wait", "w", 30, "Seconds to wait between iterations")

	monitorCmd.MarkFlagRequired("url")

	// This allows the command line flags to overwrite the configuration file values
	viper.BindPFlag("url", monitorCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("username", monitorCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", monitorCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("insecure", monitorCmd.PersistentFlags().Lookup("insecure"))
	viper.BindPFlag("count", monitorCmd.PersistentFlags().Lookup("count"))
	viper.BindPFlag("wait", monitorCmd.PersistentFlags().Lookup("wait"))
	viper.BindPFlag("pgurl", monitorCmd.PersistentFlags().Lookup("pgurl"))
	//(&Url, "url", "URL", "Cockroach Cluster API URL")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
