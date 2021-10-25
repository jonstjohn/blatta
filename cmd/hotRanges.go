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
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"math"
	"net/url"
	"sort"
	"text/tabwriter"
	"time"

	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

//type loginResponseStruct struct {
//	Session string
//}

type RangeResponse struct {
	RangeId int `json:"range_id"`
	EndKey string `json:"end_key"`
	StoreId int `json:"store_id"`
	QueriesPerSecond float32 `json:"queries_per_second"`
}

type Range struct {
	NodeId string
	RangeId int
	StoreId int
	QueriesPerSecond float32
	Database string
	TableName string
	StartPretty string
	EndPretty string
}

type RangesByNodeIdResponse map[string][]RangeResponse

type HotRangesResponse struct {
	RangesByNodeId RangesByNodeIdResponse `json:"ranges_by_node_id"`
	Next string
}

func getHotRangesResponse(apiKey string, apiUrl string) HotRangesResponse {

	hotRangeResource := "/api/v2/ranges/hot/"
	uHr, _ := url.ParseRequestURI(apiUrl)
	uHr.Path = hotRangeResource
	urlStrHr := uHr.String()

	r, _ := http.NewRequest(http.MethodGet, urlStrHr, nil) // URL-encoded payload
	r.Header.Add("X-Cockroach-API-Session", apiKey)

	client := api.HttpClient(apiUrl, Insecure)

	resp, _ := client.Do(r)
	body, _ := ioutil.ReadAll(resp.Body)

	var hotRangesResponse HotRangesResponse
	json.Unmarshal(body, &hotRangesResponse)
	return hotRangesResponse
}

func sortRangesWithNodeId(nodeRanges RangesByNodeIdResponse) []Range {
	// Convert ranges to include the node ID which we want attached to each range
	var allRanges []Range

	for nodeId, ranges := range nodeRanges {
		for _, r := range ranges {
			allRanges = append(allRanges, Range{
				NodeId: nodeId,
				RangeId: r.RangeId,
				StoreId: r.StoreId,
				QueriesPerSecond: r.QueriesPerSecond,
			})
		}
	}

	// Sort ranges from highest QPS to lowest
	sort.SliceStable(allRanges, func(i, j int) bool {
		return allRanges[i].QueriesPerSecond > allRanges[j].QueriesPerSecond
	})

	return allRanges
}

func populateAdditionalRangeInfo(allRanges []Range, pgUrl string) error {
	conn, err := pgx.Connect(context.Background(), pgUrl)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	for i, r := range allRanges {
		var tableId int
		var systemTableName string
		var indexName string
		var startPretty string
		var endPretty string
		err = conn.QueryRow(
			context.Background(),
			"select table_id, table_name, index_name, start_pretty, end_pretty from crdb_internal.ranges where range_id = $1", r.RangeId,
		).Scan(&tableId, &systemTableName, &indexName, &startPretty, &endPretty)
		if err != nil {
			return err
		}

		var databaseName string
		var tableName string

		if tableId == 0 {
			tableName = systemTableName
		} else {
			err = conn.QueryRow(
				context.Background(),
				"select database_name, name from crdb_internal.tables where table_id = $1", tableId,
			).Scan(&databaseName, &tableName)
			if err != nil {
				return err
			}
		}

		maxLength := 30
		var startPrettyStr string
		if len(startPretty) <= maxLength {
			startPrettyStr = startPretty
		} else {
			startPrettyStr = fmt.Sprintf("%s...", startPretty[0:maxLength])
		}
		var endPrettyStr string
		if len(endPretty) <= maxLength {
			endPrettyStr = endPretty
		} else {
			endPrettyStr = fmt.Sprintf("%s...", endPretty[0:maxLength])
		}

		if len(indexName) > 0 {
			tableName = fmt.Sprintf("%s@%s", tableName, indexName)
		}
		r.Database = databaseName
		r.TableName = tableName
		r.StartPretty = startPrettyStr
		r.EndPretty = endPrettyStr
		allRanges[i] = r
	}

	return nil
}

func printRanges(ranges []Range) {

	t := time.Now()
	fmt.Println(t.Format(time.RFC3339))

	w := tabwriter.NewWriter(os.Stdout, 7, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w,"Node\tRange ID\tQPS\tStore\tDB\tTable / Index\tStart Key\tEnd Key\t")

	for _, r := range ranges {
		fmt.Fprintf(w,"%s\t%d\t%.2f\t%d\t%s\t%s\t%s\t%s\t\n", r.NodeId, r.RangeId,
			r.QueriesPerSecond, r.StoreId, r.Database, r.TableName,
			r.StartPretty, r.EndPretty)
	}
	w.Flush()
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

		apiUrl := viper.GetString("url")
		username := viper.GetString("username")
		password := viper.GetString("password")
		//cacert := viper.GetString("cacert")
		pgUrl := viper.GetString("pgurl")
		insecure := viper.GetBool("insecure")

		// Login with username and password to get API key
		apiKey, err := api.Login(apiUrl, username, password, insecure)
		if err != nil {
			return err
		}

		// Iterate for "Count" iterations, or use max int if zero
		iterations := Count
		if Count == 0 {
			iterations = math.MaxInt8
		}
		for i := 1; i < iterations; i++ {

			// Get ranges by node ID from the response (page of response)
			hotRangesResponse := getHotRangesResponse(apiKey, apiUrl)

			// Sort ranges from highest QPS to lowest and add node ID, take first 10
			allRanges := sortRangesWithNodeId(hotRangesResponse.RangesByNodeId)[0:10]

			populateAdditionalRangeInfo(allRanges, pgUrl)
			printRanges(allRanges)

			time.Sleep(time.Duration(Wait) * time.Second)

		}

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
