package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	EuiUrl                   *string
	EuiConfig                *EuiConfig
	SelectedDTOFilters       *[]string
	SelectedAttributeFilters *[]string
	Es                       *Es
	Validity                 *Validity
	SearchTypes              *[]SearchType
	AddedGroupRights         *[]GroupRight
	AddedUserDbs             *[]Userdb
	UserDBConfig             *UserDBConfig
	Response                 http.Response
}

func (c *Client) SendConfig() {
	url := fmt.Sprintf("%s/eui/config", *c.EuiUrl)

	configJson, err := json.Marshal(c.EuiConfig)
	if err != nil {
		panic("Unable to parse json")
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		panic(fmt.Sprintf("Unable to send request: %s", err))
	}
	c.Response = *resp
}

// TODO: Cleanup/logging/errorhandling
func (c *Client) SendFilters() {
	client := &http.Client{}

	url := fmt.Sprintf("%s/eui/config/filters", *c.EuiUrl)
	selectedAttributeFilters := GetFilters(WIZEPASS_ATTRIBUTE_OPTIONS, *c.SelectedAttributeFilters)
	selectedDtoFilters := GetFilters(WIZEPASS_DTO_OPTIONS, *c.SelectedDTOFilters)

	for _, filter := range selectedAttributeFilters {
		configJson, err := json.Marshal(filter)
		if err != nil {
			panic("Unable to parse json")
		}

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

			printResponse("AttributeFilter", res)

		}

	}

	for _, filter := range selectedDtoFilters {
		configJson, err := json.Marshal(filter)
		if err != nil {
			panic("Unable to parse json")
		}

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
		}

	}
}

func (c *Client) SendEsConnection() {
	url := fmt.Sprintf("%s/eui/config/es", *c.EuiUrl)

	c.Es.Validity.SetDurationFromDays()

	configJson, err := json.Marshal(c.Es)
	if err != nil {
		panic("Unable to parse json")
	}

	fmt.Println(string(configJson))

	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		panic(fmt.Sprintf("Unable to send request: %s", err))
	}
	printResponse("Es", resp)
}

// TODO: Cleanup/logging/errorhandling
func (c *Client) SendSearchTypes() {
	url := fmt.Sprintf("%s/eui/config/search-types", *c.EuiUrl)

	for _, searchType := range *c.SearchTypes {
		searchJson, err := json.Marshal(searchType)
		if err != nil {
			panic("Unable to parse json")
		}

		resp, err := http.Post(url, "application/json", bytes.NewReader(searchJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("SearchTypes", resp)
	}
}

// TODO: Cleanup/logging/errorhandling
// TODO: replace with added grouprighs to be able to add userdbs
func (c *Client) SendGroupRights() {
	url := fmt.Sprintf("%s/eui/config/rights", *c.EuiUrl)

	for _, groupRight := range *c.AddedGroupRights {
		grouprightJson, err := json.Marshal(groupRight)
		if err != nil {
			panic("Unable to parse json")
		}

		resp, err := http.Post(url, "application/json", bytes.NewReader(grouprightJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("Group rights", resp)
	}
}

func (c *Client) SendUserdbConnection() {
	url := fmt.Sprintf("%s/eui/config/userdb", *c.EuiUrl)

	for _, userdb := range *c.AddedUserDbs {
		userdbJson, err := json.Marshal(userdb)
		if err != nil {
			panic("Unable to parse json")
		}

		resp, err := http.Post(url, "application/json", bytes.NewReader(userdbJson))
		if err != nil {
			panic(fmt.Sprintf("Unable to send request: %s", err))
		}
		printResponse("Userdb", resp)
	}
}

func printResponse(title string, res *http.Response) {
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println(title)
	fmt.Println("Response Status: ", res.Status)
	fmt.Println("Response Body: ", string(bodyBytes))
}
