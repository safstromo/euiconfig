package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func EuiConfigForm(client *Client, newConfig *Config) {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	euiConfigForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Value(&newConfig.EuiConfig.EsUrl).
				Title("Enter ES Url").
				Placeholder("https://api3.test.wizepass.com").
				Description("Url to Enrolment Service"),

			huh.NewInput().
				Value(&newConfig.EuiConfig.RpUrl).
				Title("Enter RP Url").
				Placeholder("https://api3.test.wizepass.com").
				Description("Url to Relying Party Service"),

			huh.NewInput().
				Value(&newConfig.EuiConfig.RpSignId).
				Title("Enter RP sign id").
				Placeholder("https://api3.test.wizepass.com").
				Description("UUID to RP for signing"),

			huh.NewConfirm().
				Title("Rp request required").
				Value(&newConfig.EuiConfig.RpRequestRequired).
				Affirmative("Yes!").
				Negative("No.").
				Description("If rp is reqired for sign"),
		),
	).WithAccessible(accessible)

	err := euiConfigForm.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	_ = spinner.New().Title("Sending Euiconfig...").Accessible(accessible).Action(client.SendConfig).Run()
}
