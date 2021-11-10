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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// debugZipCmd represents the debugZip command
var debugZipCmd = &cobra.Command{
	Use:   "debugZip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("debugZip called")
		fmt.Println(cmd.Flag("filepath").Value)
		filepath, _ := cmd.Flags().GetString("filepath") //cmd.Flag("filepath").Value

		_, err := isDebugZip(filepath)
		if err != nil {
			return err
		}

		return nil
	},
}

func isDebugZip(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.Name() == "nodes.json" {
			return true, nil
		}
	}

	return false, errors.New("not a valid debug zip - could not find nodes.json")


	return true, nil
}

func init() {
	analyzeCmd.AddCommand(debugZipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// debugZipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// debugZipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	debugZipCmd.Flags().String("filepath", "", "Debug zip file path")
	debugZipCmd.MarkFlagRequired("filepath")
}
