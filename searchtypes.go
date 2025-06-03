package main

import (
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type SearchType struct {
	Entity string `json:"entity"`
	Type   string `json:"type"`
	State  bool   `json:"state"`
}

// TODO: fix description/cleanup/refactor
func SearchTypeForm(newConfig *Config) {
	Log.Info("Starting searchtypes form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	// Reset to be able to run again
	newConfig.AddedTypes = []SearchType{}

	moreTypes := true

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

		RunForm(searchTypeForm)

		newConfig.AddedTypes = append(newConfig.AddedTypes, newSearchType)
	}

	_ = spinner.New().Title("Sending Searchtypes config...").Accessible(accessible).Action(newConfig.SendSearchTypes).Run()
}
