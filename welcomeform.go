package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func WelcomeForm(newConfig *Config) {
	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	Log.Info("Starting welcome form")
	welcome := fmt.Sprintf("Welcome to _EuiConfig™_.\n\n%s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("_!!Dont forget to disable auth before you continue!!_"))

	form := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Eui Config").
			Description(welcome),

			huh.NewInput().
				Value(&newConfig.EuiUrl).
				Title("Please enter Eui Url").
				Placeholder(newConfig.EuiUrl).
				Description("The url to the Eui to configure"),
		),
	).WithAccessible(accessible)

	RunForm(form)
}
