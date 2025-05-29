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

	fmt.Println("sel filter:", len(selectedAttributeFilters))
	fmt.Println("dto filter:", len(selectedDtoFilters))

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
