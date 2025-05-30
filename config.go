package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	EuiUrl                   string
	EuiConfig                EuiConfig
	SelectedDTOFilters       []string
	SelectedAttributeFilters []string
	Es                       Es
	AddedTypes               []SearchType
	AddedGroupRights         []GroupRight
	AddedUserDbs             []Userdb
	UserDBConfig             UserDBConfig
	Response                 http.Response
}

func (c *Config) SendConfig() {
	Log.Info("Starting send config")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config", c.EuiUrl)

	configJson := CreateJson(c.EuiConfig)

	Log.Info("Sending Eui config")
	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		Log.Errorf("Unable to send request: %s", err)
	}
	printResponse("Eui Config", resp)
}

// TODO: Cleanup/logging/errorhandling
func (c *Config) SendFilters() {
	Log.Info("Starting send filters")
	time.Sleep(1 * time.Second)
	client := &http.Client{}

	url := fmt.Sprintf("%s/eui/config/filters", c.EuiUrl)
	selectedAttributeFilters := GetFilters(WIZEPASS_ATTRIBUTE_OPTIONS, c.SelectedAttributeFilters)
	selectedDtoFilters := GetFilters(WIZEPASS_DTO_OPTIONS, c.SelectedDTOFilters)

	Log.Info("Looping through and sending filters")
	for _, filter := range selectedAttributeFilters {

		configJson := CreateJson(filter)

		resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}

		if resp.StatusCode == 400 {
			fmt.Println("Unable to create filter sending PUT")
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(configJson))
			if err != nil {
				Log.Errorf("Unable to send request: %s", err)
				panic(fmt.Sprintf("Unable to send request: %s", err))
			}

			res, _ := client.Do(req)

			printResponse("AttributeFilter", res)
		} else {
			printResponse("AttributeFilter", resp)
		}

	}

	for _, filter := range selectedDtoFilters {

		configJson := CreateJson(filter)

		resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}

		if resp.StatusCode == 400 {
			fmt.Println("Unable to create filter sending PUT")
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(configJson))
			if err != nil {
				panic(fmt.Sprintf("Unable to send request: %s", err))
			}

			res, _ := client.Do(req)

			printResponse("DtoFilter", res)
		} else {
			printResponse("DtoFilter", resp)
		}

	}
}

func (c *Config) SendEsConnection() {
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/es", c.EuiUrl)

	c.Es.Validity.SetDurationFromDays()

	configJson := CreateJson(c.Es)

	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		panic(fmt.Sprintf("Unable to send request: %s", err))
	}
	printResponse("Es Connection", resp)
}

// TODO: Cleanup/logging/errorhandling
func (c *Config) SendSearchTypes() {
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/search-types", c.EuiUrl)

	for _, searchType := range c.AddedTypes {
		searchJson := CreateJson(searchType)

		resp, err := http.Post(url, "application/json", bytes.NewReader(searchJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("SearchTypes", resp)
	}
}

// TODO: Cleanup/logging/errorhandling
// TODO: replace with added grouprighs to be able to add userdbs
func (c *Config) SendGroupRights() {
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/rights", c.EuiUrl)

	for _, groupRight := range c.AddedGroupRights {
		grouprightJson := CreateJson(groupRight)

		resp, err := http.Post(url, "application/json", bytes.NewReader(grouprightJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("GroupRights", resp)
	}
}

func (c *Config) SendUserdbConnection() {
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/userdb", c.EuiUrl)

	for _, userdb := range c.AddedUserDbs {
		userdbJson := CreateJson(userdb)

		resp, err := http.Post(url, "application/json", bytes.NewReader(userdbJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("Userdb connection", resp)
	}
}

func (c *Config) SendUserDbConfig() {
	time.Sleep(1 * time.Second)
	client := &http.Client{}

	url := fmt.Sprintf("%s/eui/config/userdb/config", c.EuiUrl)

	configJson, err := json.Marshal(c.UserDBConfig)
	if err != nil {
		panic("Unable to parse json")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(configJson))
	if err != nil {
		panic(fmt.Sprintf("Unable to send request: %s", err))
	}
	req.Header.Add("Content-Type", "application/json")

	userdbUrl := (c.AddedUserDbs)[0].Url

	query := req.URL.Query()
	query.Add("url", userdbUrl)
	req.URL.RawQuery = query.Encode()

	res, _ := client.Do(req)

	printResponse("Userdb config", res)
}

func ReadBody(res *http.Response) string {
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "No body"
	}

	return string(bodyBytes)
}

func printResponse(title string, response *http.Response) {
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}

	fmt.Fprintf(&sb, "%s\n\n", lipgloss.NewStyle().Bold(true).Render(title+" response: "))

	fmt.Fprintf(&sb,
		"%s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Render(response.Status),
		lipgloss.NewStyle().Bold(true).Render("Body: "),
		keyword(ReadBody(response)),
	)

	if response.StatusCode == 200 || response.StatusCode == 201 {
		fmt.Fprint(&sb, "\n\nGreat success!")
	} else {
		fmt.Fprint(&sb, "\n\nSomething went wrong")
	}

	fmt.Println(
		lipgloss.NewStyle().
			Width(80).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)
}

func CreateJson(object any) []byte {
	Log.Info("Creating json")
	json, err := json.Marshal(object)
	if err != nil {
		Log.Error("Unable to parse json")
		panic("Unable to parse json")
	}
	return json
}
