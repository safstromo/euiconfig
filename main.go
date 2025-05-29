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
)

type EuiConfig struct {
	EsUrl             string   `json:"es_url"`
	RpUrl             string   `json:"rp_url"`
	RpSignId          string   `json:"rp_sign_id"`
	RpRequestRequired bool     `json:"rp_request_required"`
	RevokeComments    []string `json:"revoke_comments"`
}

type Config struct {
	EuiUrl                   string
	EuiConfig                EuiConfig
	SelectedDTOFilters       []string
	SelectedAttributeFilters []string
	Es                       Es
	AddedTypes               []SearchType
	AddedGroupRights         []GroupRight
	AddedUserDbs             []Userdb
	Response                 http.Response
}

// TODO:show request status  for user
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
		Es: Es{
			State:    true,
			UniqueID: "a0759625-f432-41f6-9dca-0c42e51aa1d5",
			Validity: Validity{
				UseDuration: true,
			},
		},
	}

	// dtoOptions := ConvertFiltersToHuhOptions(WIZEPASS_DTO_OPTIONS)
	// attributeOptions := ConvertFiltersToHuhOptions(WIZEPASS_ATTRIBUTE_OPTIONS)

	client := Client{
		EuiUrl:                   &newConfig.EuiUrl,
		EuiConfig:                &newConfig.EuiConfig,
		SelectedDTOFilters:       &newConfig.SelectedDTOFilters,
		SelectedAttributeFilters: &newConfig.SelectedAttributeFilters,
		Es:                       &newConfig.Es,
		Validity:                 &newConfig.Es.Validity,
		SearchTypes:              &newConfig.AddedTypes,
		AddedGroupRights:         &newConfig.AddedGroupRights,
		AddedUserDbs:             &newConfig.AddedUserDbs,
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

		// // EuiConfig
		// huh.NewGroup(
		// 	huh.NewInput().
		// 		Value(&newConfig.EuiConfig.EsUrl).
		// 		Title("Enter ES Url").
		// 		Placeholder("https://api3.test.wizepass.com").
		// 		Description("Url to Enrolment Service"),
		//
		// 	huh.NewInput().
		// 		Value(&newConfig.EuiConfig.RpUrl).
		// 		Title("Enter RP Url").
		// 		Placeholder("https://api3.test.wizepass.com").
		// 		Description("Url to Relying Party Service"),
		//
		// 	huh.NewInput().
		// 		Value(&newConfig.EuiConfig.RpSignId).
		// 		Title("Enter RP sign id").
		// 		Placeholder("https://api3.test.wizepass.com").
		// 		Description("UUID to RP for signing"),
		//
		// 	huh.NewConfirm().
		// 		Title("Rp request required").
		// 		Value(&newConfig.EuiConfig.RpRequestRequired).
		// 		Affirmative("Yes!").
		// 		Negative("No.").
		// 		Description("If rp is reqired for sign"),
		// ),
		//
		// // Filters
		// huh.NewGroup(
		// 	huh.NewMultiSelect[string]().
		// 		Options(dtoOptions...).
		// 		Title("Wizepass DTO Filters").
		// 		Value(&newConfig.SelectedDTOFilters),
		// ),
		// huh.NewGroup(
		// 	huh.NewMultiSelect[string]().
		// 		Options(attributeOptions...).
		// 		Title("Wizepass Attribute Filters").
		// 		Value(&newConfig.SelectedAttributeFilters),
		// ),
		//
		// // Es connection
		// // TODO: tags/ Default confige
		// huh.NewGroup(
		// 	huh.NewNote().
		// 		Title("Create ES connection"),
		// 	huh.NewInput().
		// 		Value(&newConfig.Es.Name).
		// 		Title("Name").
		// 		Placeholder("Default Es"),
		// 	huh.NewInput().
		// 		Value(&newConfig.Es.UniqueID).
		// 		Title("UUID").
		// 		Placeholder("a0759625-f432-41f6-9dca-0c42e51aa1d5").
		// 		Description("UUID of Enrolment Service"),
		// 	huh.NewInput().
		// 		Value(&newConfig.Es.URL).
		// 		Title("Client Url").
		// 		Placeholder("client.wizepass.com").
		// 		Description("Url for client, dont include https://"),
		// 	huh.NewInput().
		// 		Value(&newConfig.Es.Validity.DisplayName).
		// 		Title("Validity name").
		// 		Placeholder("One year"),
		// 	huh.NewInput().
		// 		Value(&newConfig.Es.Validity.Description).
		// 		Title("Validity Description").
		// 		Placeholder("One year duration"),
		// 	huh.NewInput().
		// 		Value(&newConfig.Es.Validity.DurationString).
		// 		Title("Validity Duration").
		// 		Placeholder("365").
		// 		Description("Duration in days"),
		// ),
	).WithAccessible(accessible)

	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// _ = spinner.New().Title("Sending Euiconfig...").Accessible(accessible).Action(client.SendConfig).Run()
	// _ = spinner.New().Title("Sending filters...").Accessible(accessible).Action(client.SendFilters).Run()
	// _ = spinner.New().Title("Sending Es connection config...").Accessible(accessible).Action(client.SendEsConnection).Run()

	moreTypes := true

	// TODO: fix description/cleanup/refactor
	for moreTypes {
		newSearchType := SearchType{
			State:  true,
			Entity: "wizepass",
		}

		searchTypeForm := huh.NewForm(huh.NewGroup(
			huh.NewNote().
				Title("Add SearchTypes").
				Description("EUi seachtypes.\n"),

			huh.NewInput().
				Value(&newSearchType.Type).
				Title("Please enter type").
				Placeholder("user_id").
				Description("Fields to search for"),

			huh.NewConfirm().
				Title("Add more searchtypes?").
				Value(&moreTypes).
				Affirmative("Yes!").
				Negative("No.").
				Description("Continue to add more searchTypes"),
		)).WithAccessible(accessible)

		err := searchTypeForm.Run()
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}
		newConfig.AddedTypes = append(newConfig.AddedTypes, newSearchType)
	}

	_ = spinner.New().Title("Sending Searchtypes config...").Accessible(accessible).Action(client.SendSearchTypes).Run()

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
