package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	// xstrings "github.com/charmbracelet/x/exp/strings"
)

type EuiConfig struct {
	EsUrl             string   `json:"es_url"`
	RpUrl             string   `json:"rp_url"`
	RpSignId          string   `json:"rp_sign_id"`
	RpRequestRequired bool     `json:"rp_request_required"`
	RevokeComments    []string `json:"revoke_comments"`
}

type Config struct {
	EuiUrl    string
	EuiConfig EuiConfig
	Response  http.Response
}

func main() {
	newConfig := Config{
		EuiUrl: "http://localhost:8080",
		EuiConfig: EuiConfig{
			EsUrl:             "https://api3.test.wizepass.com",
			RpUrl:             "https://api3.test.wizepass.com",
			RpSignId:          "a0759625-f432-41f6-9dca-0c42e51aa1d5",
			RpRequestRequired: false,
			RevokeComments: []string{
				"unspecified",
				"key_compromise",
				"affiliation_changed",
				"superseded",
				"cessation_of_operation",
				"privilege_withdrawn",
			},
		},
	}

	client := Client{
		EuiUrl:    &newConfig.EuiUrl,
		EuiConfig: &newConfig.EuiConfig,
	}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Eui Config").
			Description("Welcome to _EuiConfig™_.\n"),

			huh.NewInput().
				Value(&newConfig.EuiUrl).
				Title("Please enter Eui Url").
				Placeholder(newConfig.EuiUrl).
				Description("The url to the Eui to configure"),
		),

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

	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	_ = spinner.New().Title("Sending config...").Accessible(accessible).Action(client.SendConfig).Run()

	{
		var sb strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}

		fmt.Fprintf(&sb, "%s\n\n", lipgloss.NewStyle().Bold(true).Render("CONFIG"))

		fmt.Fprintf(&sb,
			"%s%s\n",
			lipgloss.NewStyle().Bold(true).Render("EUI Url: "),
			keyword(newConfig.EuiUrl),
		)

		fmt.Fprintf(&sb,
			"%s%s\n",
			lipgloss.NewStyle().Bold(true).Render("Es Url: "),
			keyword(newConfig.EuiConfig.EsUrl),
		)
		fmt.Fprintf(&sb,
			"%s%s\n",
			lipgloss.NewStyle().Bold(true).Render("Rp Url: "),
			keyword(newConfig.EuiConfig.RpUrl),
		)

		if newConfig.EuiConfig.RpRequestRequired {
			fmt.Fprintf(&sb,
				"%s%s\n",
				lipgloss.NewStyle().Bold(true).Render("Rp request Required: "),
				keyword("true"),
			)
			fmt.Fprintf(&sb,
				"%s%s\n",
				lipgloss.NewStyle().Bold(true).Render("Rp Sign ID: "),
				keyword(newConfig.EuiConfig.RpSignId),
			)
		}

		fmt.Fprintf(&sb,
			"%s\n",
			lipgloss.NewStyle().Bold(true).Render("Revoke comments:"),
		)
		for _, v := range newConfig.EuiConfig.RevokeComments {
			fmt.Fprintf(&sb,
				"%s\n",
				keyword(v),
			)
		}

		fmt.Fprintf(&sb, "\n\n%s", client.Response.Status)
		fmt.Fprint(&sb, "\n\nThanks for using EuiConfig!")

		fmt.Println(
			lipgloss.NewStyle().
				Width(80).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}
