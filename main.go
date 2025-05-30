package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Log *log.Logger

func init() {
	logFile, err := os.OpenFile("euiConfig.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	Log = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	})

	Log.SetOutput(logFile)
	Log.Info("Logger initialized.")
}

// TODO: connect userdb with groupright
func main() {
	Log.Info("Creating default config")
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
		UserDBConfig: NewDefaultUserDBConfig(),
	}

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

	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	EuiConfigForm(&newConfig)
	FiltersForm(&newConfig)
	EsConnectionForm(&newConfig)
	SearchTypeForm(&newConfig)
	GroupRightsForm(&newConfig)
	UserDbConnectionForm(&newConfig)
	UserDbConfigForm(&newConfig)

	{
		Log.Info("Printing goodbye")
		var sb strings.Builder

		fmt.Fprintf(&sb, "\n%s\n", lipgloss.NewStyle().Bold(true).Render("Thanks for using EuiConfig"))

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
