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
	"sort"
	"strconv"
	"strings"

	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

type loginResponseStruct struct {
	Session string
}

type Range struct {
	RangeId int `json:"range_id"`
	EndKey string `json:"end_key"`
	StoreId int `json:"store_id"`
	QueriesPerSecond float32 `json:"queries_per_second"`
}

type RangesByNodeId map[string][]Range

type HotRangesResponse struct {
	RangesByNodeId RangesByNodeId `json:"ranges_by_node_id"`
	Next string
}

func HttpClient(apiUrl string, insecure bool) *http.Client {
	customTransport := &(*http.DefaultTransport.(*http.Transport)) // make shallow copy
	if insecure {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &http.Client{Transport: customTransport}
}

func Login(apiUrl string, username string, password string, insecure bool) (string, error) {

	resource := "/api/v2/login/"
	data := url.Values{}
	data.Set("username", Username)
	data.Set("password", Password)

	u, error := url.ParseRequestURI(apiUrl)

	if error != nil {
		return "", error
	}

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

	return apiKey, nil
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
	RunE: func(cmd *cobra.Command, args []string) error {

		apiKey, err := Login(ApiUrl, Username, Password, Insecure)

		if err != nil {
			return err
		}

		hotRangeResource := "/api/v2/ranges/hot/"
		uHr, _ := url.ParseRequestURI(ApiUrl)
		uHr.Path = hotRangeResource
		urlStrHr := uHr.String()

		r, _ := http.NewRequest(http.MethodGet, urlStrHr, nil) // URL-encoded payload
		r.Header.Add("X-Cockroach-API-Session", apiKey)

		client := HttpClient(ApiUrl, Insecure)

		resp, _ := client.Do(r)
		body, _ := ioutil.ReadAll(resp.Body)
		//bodyString := string(body)

		var hotRangesResponse HotRangesResponse
		json.Unmarshal(body, &hotRangesResponse)
		//decoder := json.NewDecoder(bodyString)
		fmt.Printf("Hot ranges response: %+v", hotRangesResponse)

		nodeRanges := hotRangesResponse.RangesByNodeId
		fmt.Printf("Node ranges: %+v", nodeRanges)

		var allRanges []Range

		for nodeId, ranges := range nodeRanges {
			fmt.Println(nodeId)
			for _, r := range ranges {
				allRanges = append(allRanges, r)
			}
		}

		sort.SliceStable(allRanges, func(i, j int) bool {
			return allRanges[i].QueriesPerSecond > allRanges[j].QueriesPerSecond
		})

		fmt.Printf("%+v", allRanges[0:10])

		return nil
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
