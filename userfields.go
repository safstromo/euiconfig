package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	// "github.com/charmbracelet/huh/spinner"
)

type UserFilter struct {
	Entity string `json:"entity"` // user_dto
	Field  string `json:"field"`  // name
	State  bool   `json:"state"`
}

// TODO:no endpoint
func UserDbFilterForm(newConfig *Config) {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	moreFields := true

	for moreFields {
		newFilter := UserFilter{
			State:  true,
			Entity: "user_dto",
		}

		searchTypeForm := huh.NewForm(huh.NewGroup(
			huh.NewNote().
				Title("Add Userdb filter").
				Description("Userdb User fields.\n"),

			huh.NewInput().
				Value(&newFilter.Field).
				Title("Please enter field").
				Placeholder("user_id").
				Description("Fields to filter."),

			huh.NewConfirm().
				Title("Add more filters?").
				Value(&moreFields).
				Affirmative("Yes!").
				Negative("No.").
				Description("Continue to add more filters"),
		)).WithAccessible(accessible)

		err := searchTypeForm.Run()
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}
		// newConfig.AddedTypes = append(newConfig.AddedTypes, newFilter)
	}

	// _ = spinner.New().Title("Sending Searchtypes config...").Accessible(accessible).Action(client.SendUserDbFilters).Run()
}
