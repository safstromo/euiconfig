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
	AddedGroupRights         []GroupRightPost
	GroupRightsPut           []GroupRightsPut
	AddedUserDbs             []Userdb
	UserDBConfig             UserDBConfig
	Response                 http.Response
}

func (c *Config) SendConfig() {
	Log.Info("Starting send config")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config", c.EuiUrl)

	configJson := CreateJson(c.EuiConfig)

	Log.Info("Sending request")
	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		Log.Errorf("Unable to send request: %s", err)
	}
	printResponse("Eui Config", resp)
}

// TODO: Cleanup
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

		Log.Info("Sending request")
		resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}

		if resp.StatusCode == 400 {
			Log.Error("Unable to create filter sending PUT")
			Log.Info("Sending request")
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(configJson))
			if err != nil {
				Log.Errorf("Unable to send request: %s", err)
				panic(fmt.Sprintf("Unable to send request: %s", err))
			}

			Log.Info("Sending request")
			res, _ := client.Do(req)

			printResponse("AttributeFilter", res)
		} else {
			printResponse("AttributeFilter", resp)
		}

	}

	for _, filter := range selectedDtoFilters {
		configJson := CreateJson(filter)

		Log.Info("Sending request")
		resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}

		if resp.StatusCode == 400 {
			fmt.Println("Unable to create filter sending PUT")
			Log.Info("Sending request")
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(configJson))
			if err != nil {
				Log.Errorf("Unable to send request: %s", err)
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
	Log.Info("Starting send es connection")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/es", c.EuiUrl)

	c.Es.Validity.SetDurationFromDays()

	configJson := CreateJson(c.Es)

	Log.Info("Sending request")
	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		Log.Errorf("Unable to send request: %s", err)
		panic(fmt.Sprintf("Unable to send request: %s", err))
	}
	printResponse("Es Connection", resp)
}

// TODO: Cleanup
func (c *Config) SendSearchTypes() {
	Log.Info("Starting send searchtypes")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/search-types", c.EuiUrl)

	for _, searchType := range c.AddedTypes {
		searchJson := CreateJson(searchType)

		Log.Info("Sending request")
		resp, err := http.Post(url, "application/json", bytes.NewReader(searchJson))
		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("SearchTypes", resp)
	}
}

// TODO: Cleanup
// TODO: replace with added grouprighs to be able to add userdbs
func (c *Config) SendGroupRights() {
	Log.Info("Starting send grouprights")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/rights", c.EuiUrl)

	for _, groupRight := range c.AddedGroupRights {
		grouprightJson := CreateJson(groupRight)

		Log.Info("Sending request")
		resp, err := http.Post(url, "application/json", bytes.NewReader(grouprightJson))
		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			Log.Errorf("Error reading response body: %v\n", err)
			printResponse("Userdb connection", resp)
			return
		}

		var addedGroupRight GroupRightsPut

		Log.Info("Replacing added groupright")
		err = json.Unmarshal(bodyBytes, &addedGroupRight)
		if err != nil {
			Log.Errorf("Error unmarshaling JSON: %v\n", err)
			printResponse("Userdb connection", resp)
		}

		c.GroupRightsPut = append(c.GroupRightsPut, addedGroupRight)
		printResponse("GroupRights", resp)
	}
}

func (c *Config) SendConnectGroupRights() {
	Log.Info("Starting send connect grouprights")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/rights", c.EuiUrl)
	client := &http.Client{}

	for _, groupRight := range c.GroupRightsPut {
		grouprightJson := CreateJson(groupRight)

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(grouprightJson))
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}

		Log.Info("Sending request")
		res, _ := client.Do(req)

		printResponse("GroupRights", res)
	}
}

func (c *Config) SendUserdbConnection() {
	Log.Info("Starting send userdb connection")
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/eui/config/userdb", c.EuiUrl)

	for _, userdb := range c.AddedUserDbs {
		userdbJson := CreateJson(userdb)

		Log.Info("Sending request")
		resp, err := http.Post(url, "application/json", bytes.NewReader(userdbJson))
		if err != nil {
			Log.Errorf("Unable to send request: %s", err)
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			Log.Errorf("Error reading response body: %v\n", err)
			printResponse("Userdb connection", resp)
			return
		}

		var userdb Userdb
		err = json.Unmarshal(bodyBytes, &userdb)
		if err != nil {
			Log.Errorf("Error unmarshaling JSON: %v\n", err)
			printResponse("Userdb connection", resp)
		}

		c.AddedUserDbs[0] = userdb
		printResponse("Userdb connection", resp)
	}
}

func (c *Config) SendUserDbConfig() {
	Log.Info("Starting send userdb config")
	time.Sleep(1 * time.Second)
	client := &http.Client{}

	url := fmt.Sprintf("%s/eui/config/userdb/config", c.EuiUrl)

	configJson := CreateJson(c.UserDBConfig)

	Log.Info("Sending request")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(configJson))
	if err != nil {
		Log.Errorf("Unable to send request: %s", err)
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
	Log.Info("Reading body")
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "No body"
	}

	return string(bodyBytes)
}

func printResponse(title string, response *http.Response) {
	Log.Infof("Creating response print for response: \n%v", response)
	Log.Infof("Body: \n%s", string(ReadBody(response)))
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
