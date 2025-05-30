package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type GroupRightPost struct {
	Name             string   `json:"name"`
	Roles            []string `json:"roles"`
	AllowedUserDbIds []int    `json:"allowed_userdb_ids"`
}

type GroupRightsPut struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Roles          []string `json:"roles"`
	AllowedUserDbs []Userdb `json:"allowed_userdbs"`
}

// TODO: fix description/cleanup/refactor
func GroupRightsForm(newConfig *Config) {
	Log.Info("Starting group rights form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	moreTypes := true
	var claims string

	for moreTypes {
		newGroupRight := GroupRightPost{
			AllowedUserDbIds: []int{},
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

	_ = spinner.New().Title("Sending Group right config...").Accessible(accessible).Action(newConfig.SendGroupRights).Run()
}

func ConnectGroupRight(newConfig *Config) {
	Log.Info("Starting connectGroupRight form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	for _, groupRight := range newConfig.AddedGroupRights {
		add := false
		userdb := newConfig.AddedUserDbs[0]

		updateGroupRightForm := huh.NewForm(huh.NewGroup(
			huh.NewNote().
				Title("Add Userdb to groupRight?"),

			huh.NewConfirm().
				Title(groupRight.Name).
				Value(&add).
				Affirmative("Yes!").
				Negative("No."),
		)).WithAccessible(accessible)

		err := updateGroupRightForm.Run()
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		if add {
			Log.Info("Attempting to add userdb to GroupRight")
			for i := range newConfig.GroupRightsPut {
				grPut := &newConfig.GroupRightsPut[i]
				if grPut.Name == groupRight.Name {
					grPut.AllowedUserDbs = append(grPut.AllowedUserDbs, userdb)
					Log.Infof("Added userdb %s to groupRight %s.", userdb.Name, grPut.Name)
				}
				Log.Infof("%v", *grPut)
			}
		} else {
			Log.Infof("User chose not to add userdb to %s.", groupRight.Name)
		}
	}
}
