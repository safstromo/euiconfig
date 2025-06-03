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

type ConfigForm int

const (
	Wizard ConfigForm = iota
	EuiConfigs
	EuiFilters
	EsConnection
	SearchTypes
	GroupRights
	UserDbConnectin
	ConnectGroupRights
	UserDbConfig
	UserDbFilter
)

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

	var configForm ConfigForm
	for {
		selectForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[ConfigForm]().
					Title("What do you want to configure?").
					Options(
						huh.NewOption("Installation Wizard", Wizard),
						huh.NewOption("Eui Configuration", EuiConfigs),
						huh.NewOption("Eui Filters", EuiFilters),
						huh.NewOption("Es Connection", EsConnection),
						huh.NewOption("Searchtypes", SearchTypes),
						huh.NewOption("Group rights", GroupRights),
						huh.NewOption("Userdb Connection", UserDbConnectin),
						huh.NewOption("Connect Group rights", ConnectGroupRights),
						huh.NewOption("Userdb Configuration", UserDbConfig),
						// huh.NewOption("Userdb Filters", UserDbFilter),
					).
					Value(&configForm),
			))
		RunForm(selectForm)

		// TODO: Add ACL form, Fix userdbfilters
		switch configForm {

		case Wizard:
			WelcomeForm(&newConfig)
			EuiConfigForm(&newConfig)
			FiltersForm(&newConfig)
			EsConnectionForm(&newConfig)
			SearchTypeForm(&newConfig)
			GroupRightsForm(&newConfig)
			UserDbConnectionForm(&newConfig)
			ConnectGroupRight(&newConfig)
			UserDbConfigForm(&newConfig)
			// UserDbFilterForm(&newConfig)
		case EuiConfigs:
			EuiConfigForm(&newConfig)
		case EuiFilters:
			FiltersForm(&newConfig)
		case EsConnection:
			EsConnectionForm(&newConfig)
		case SearchTypes:
			SearchTypeForm(&newConfig)
		case GroupRights:
			GroupRightsForm(&newConfig)
		case UserDbConnectin:
			UserDbConnectionForm(&newConfig)
		case ConnectGroupRights:
			ConnectGroupRight(&newConfig)
		case UserDbConfig:
			UserDbConfigForm(&newConfig)
			// case UserDbFilter:
			// UserDbFilterForm(&newConfig)
		}

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
