package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Log *log.Logger

func CreateLogger() (*log.Logger, *os.File) {
	logFile, err := os.OpenFile("euiConfig.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	})

	logger.SetOutput(logFile)
	Log = logger
	Log.Info("Logger initialized.")
	return logger, logFile
}

// TODO: connect userdb with groupright
func main() {
	logger, logFile := CreateLogger()
	defer logFile.Close()

	logger.Info("Creating default config")
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

	WelcomeForm(&newConfig)
	EuiConfigForm(&newConfig)
	FiltersForm(&newConfig)
	EsConnectionForm(&newConfig)
	SearchTypeForm(&newConfig)
	GroupRightsForm(&newConfig)
	UserDbConnectionForm(&newConfig)
	ConnectGroupRight(&newConfig)
	UserDbConfigForm(&newConfig)
	UserDbFilterForm(&newConfig)

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

func RunForm(form *huh.Form) {
	err := form.Run()
	if err != nil {
		{
			Log.Info("Printing goodbye")
			var sb strings.Builder

			fmt.Fprintf(&sb, "\n%s %s\n", lipgloss.NewStyle().Bold(true).Render("Uh oh:"), err)

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
		os.Exit(1)
	}
}
