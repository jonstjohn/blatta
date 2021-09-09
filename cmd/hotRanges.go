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
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

type loginResponseStruct struct {
	Session string
}

func HttpClient(apiUrl string, insecure bool) *http.Client {
	customTransport := &(*http.DefaultTransport.(*http.Transport)) // make shallow copy
	if insecure {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &http.Client{Transport: customTransport}
}

func Login(apiUrl string, username string, password string, insecure bool) string {

	resource := "/api/v2/login/"
	data := url.Values{}
	data.Set("username", Username)
	data.Set("password", Password)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := HttpClient(ApiUrl, Insecure)

	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)

	if err != nil {
		panic(err)
	}
	var apiKey string

	if resp.StatusCode == http.StatusOK {

		decoder := json.NewDecoder(resp.Body)
		var t loginResponseStruct
		err := decoder.Decode(&t)

		if err != nil {
			panic(err)
		}

		apiKey = t.Session
	}

	return apiKey
}

// hotRangesCmd represents the hotRanges command
var hotRangesCmd = &cobra.Command{
	Use:   "hotRanges",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		apiKey := Login(ApiUrl, Username, Password, Insecure)

		hotRangeResource := "/api/v2/ranges/hot/"
		uHr, _ := url.ParseRequestURI(ApiUrl)
		uHr.Path = hotRangeResource
		urlStrHr := uHr.String()

		r, _ := http.NewRequest(http.MethodGet, urlStrHr, nil) // URL-encoded payload
		r.Header.Add("X-Cockroach-API-Session", apiKey)

		/*
		customTransport := &(*http.DefaultTransport.(*http.Transport)) // make shallow copy
		if Insecure {
			customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client := &http.Client{Transport: customTransport}
		 */

		client := HttpClient(ApiUrl, Insecure)

		resp, _ := client.Do(r)
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		fmt.Println(bodyString)
	},
}

func init() {
	monitorCmd.AddCommand(hotRangesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hotRangesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hotRangesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
