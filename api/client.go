package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type loginResponseStruct struct {
	Session string
}

type NodeMetrics struct {
	SysCpuUserPercent float32 `json:"sys.cpu.user.percent"`
}

type Node struct {
	NodeId int `json:"node_id"`
	Metrics NodeMetrics `json:"metrics"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
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
	data.Set("username", username)
	data.Set("password", password)

	u, error := url.ParseRequestURI(apiUrl)

	if error != nil {
		return "", error
	}

	u.Path = resource
	urlStr := u.String()

	client := HttpClient(apiUrl, insecure)

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
	} else {
		return apiKey, errors.New(fmt.Sprintf("Login failed with result: %v", resp.StatusCode))
	}

	return apiKey, nil
}

func GetNodes(apiUrl string, apiKey string, insecure bool) []Node {
	resource := "/api/v2/nodes/"
	uHr, _ := url.ParseRequestURI(apiUrl)
	uHr.Path = resource
	urlStrHr := uHr.String()

	r, _ := http.NewRequest(http.MethodGet, urlStrHr, nil) // URL-encoded payload
	r.Header.Add("X-Cockroach-API-Session", apiKey)

	client := HttpClient(apiUrl, insecure)

	resp, _ := client.Do(r)
	body, _ := ioutil.ReadAll(resp.Body)

	var response NodesResponse
	json.Unmarshal(body, &response)
	fmt.Println(response)
	return response.Nodes
}