package main

import (
	// "errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
}

func main() {
	newConfig := Config{}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Eui Config").
			Description("Welcome to _EuiConfig™_.\n"),

			huh.NewInput().
				Value(&newConfig.EuiUrl).
				Title("Please enter Eui Url").
				Placeholder("http://localhost.com:8080").
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

	sendConfigRequest := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Sending config...").Accessible(accessible).Action(sendConfigRequest).Run()

	// Print order summary.
	{
		var sb strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}
		fmt.Fprintf(&sb,
			"%s\n\n%s ",
			lipgloss.NewStyle().Bold(true).Render("CONFIG"),
			keyword(newConfig.EuiUrl),
		)

		fmt.Fprint(&sb, "\n\nThanks for using EuiConfig!")

		fmt.Println(
			lipgloss.NewStyle().
				Width(40).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}
