package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type GroupRight struct {
	Name             string   `json:"name"`
	Roles            []string `json:"roles"`
	AllowedUserDBIDs []int    `json:"allowed_userdb_ids"`
}

// TODO: fix description/cleanup/refactor
func GroupRightsForm(client *Client, newConfig *Config) {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	moreTypes := true
	var claims string

	for moreTypes {
		newGroupRight := GroupRight{
			AllowedUserDBIDs: []int{},
		}

		groupRightForm := huh.NewForm(huh.NewGroup(
			huh.NewNote().
				Title("Add Group rights").
				Description("EUI Grouprights.\n"),

			huh.NewInput().
				Value(&newGroupRight.Name).
				Title("Please enter name").
				Placeholder("Super admin"),

			huh.NewInput().
				Value(&claims).
				Title("Enter roles").
				Placeholder("super-admin,admin").
				Description("Comma sepparated list of claims"),

			huh.NewConfirm().
				Title("Add more group rights?").
				Value(&moreTypes).
				Affirmative("Yes!").
				Negative("No.").
				Description("Continue to add more group rights"),
		)).WithAccessible(accessible)

		err := groupRightForm.Run()
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		newGroupRight.Roles = strings.Split(claims, ",")

		newConfig.AddedGroupRights = append(newConfig.AddedGroupRights, newGroupRight)
	}

	_ = spinner.New().Title("Sending Group right config...").Accessible(accessible).Action(client.SendGroupRights).Run()
}
