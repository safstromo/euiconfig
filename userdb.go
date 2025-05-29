package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type Userdb struct {
	Id   int    `json:"-"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func UserDbConnectionForm(client *Client, newConfig *Config) {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	newUserdb := Userdb{}

	userdbForm := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Userdb connection"),

			huh.NewInput().
				Value(&newUserdb.Name).
				Title("Please enter Name").
				Placeholder("Defaut Userdb"),
			huh.NewInput().
				Value(&newUserdb.Url).
				Title("Please enter Userdb Url").
				Placeholder("http://localhost:8081").
				Description("The url to the Userdb to configure"),
		),
	).WithAccessible(accessible)

	err := userdbForm.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// TODO: null?????
	fmt.Println(newUserdb.Url)

	newConfig.AddedUserDbs = append(newConfig.AddedUserDbs, newUserdb)

	_ = spinner.New().Title("Sending Userdb connection config...").Accessible(accessible).Action(client.SendUserdbConnection).Run()
}
