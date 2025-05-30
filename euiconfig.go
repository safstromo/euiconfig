package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
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
	time.Sleep(2 * time.Second)
}

func PrintEuiConfigResponse(response *http.Response) {
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}

	fmt.Fprintf(&sb, "%s\n\n", lipgloss.NewStyle().Bold(true).Render("EuiConfig response:"))

	fmt.Fprintf(&sb,
		"%s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Render(response.Status),
		lipgloss.NewStyle().Bold(true).Render("Body: "),
		keyword(ReadBody(response)),
	)

	if response.StatusCode == 200 {
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
