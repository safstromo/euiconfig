package main

import (
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type EuiConfig struct {
	EsUrl             string   `json:"es_url"`
	RpUrl             string   `json:"rp_url"`
	RpSignId          string   `json:"rp_sign_id"`
	RpRequestRequired bool     `json:"rp_request_required"`
	RevokeComments    []string `json:"revoke_comments"`
}

func EuiConfigForm(newConfig *Config) {
	Log.Info("Starting Eui config form")
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

	RunForm(euiConfigForm)

	_ = spinner.New().Title("Sending Euiconfig...").Accessible(accessible).Action(newConfig.SendConfig).Run()
}
