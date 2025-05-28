package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	EuiUrl    *string
	EuiConfig *EuiConfig
	Response  http.Response
}

func (c *Client) SendConfig() {
	configJson, err := json.Marshal(c.EuiConfig)
	if err != nil {
		panic("Unable to parse json")
	}
	url := fmt.Sprintf("%s/eui/config", *c.EuiUrl)

	resp, err := http.Post(url, "application/json", bytes.NewReader(configJson))
	if err != nil {
		panic(fmt.Sprintf("Unable to send request: %s", err))
	}
	c.Response = *resp
}
