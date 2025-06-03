package main

import (
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type Es struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	Validity   Validity `json:"validity"`
	InstanceID int      `json:"instance_id"`
	UniqueID   string   `json:"unique_id"`
	State      bool     `json:"state"`
}

type Validity struct {
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	UseDuration    bool   `json:"use_duration"`
	Duration       int64  `json:"duration"`
	DurationString string `json:"-"`
}

func (v *Validity) SetDurationFromDays() {
	Log.Info("Converting duration")
	daysInt, err := strconv.Atoi(v.DurationString)
	if err != nil {
		Log.Error("Unable to convert duration string to int: ", v.DurationString)
	}

	duration := time.Duration(daysInt) * 24 * time.Hour
	v.Duration = duration.Milliseconds()
}

// Es Connection
func EsConnectionForm(newConfig *Config) {
	Log.Info("Starting es connection form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	EsForm := huh.NewForm(
		// TODO: tags/ Default config
		huh.NewGroup(
			huh.NewNote().
				Title("Create ES connection"),
			huh.NewInput().
				Value(&newConfig.Es.Name).
				Title("Name").
				Placeholder("Default Es"),
			huh.NewInput().
				Value(&newConfig.Es.UniqueID).
				Title("UUID").
				Placeholder("a0759625-f432-41f6-9dca-0c42e51aa1d5").
				Description("UUID of Enrolment Service"),
			huh.NewInput().
				Value(&newConfig.Es.URL).
				Title("Client Url").
				Placeholder("client.wizepass.com").
				Description("Url for client, dont include https://"),
			huh.NewInput().
				Value(&newConfig.Es.Validity.DisplayName).
				Title("Validity name").
				Placeholder("One year"),
			huh.NewInput().
				Value(&newConfig.Es.Validity.Description).
				Title("Validity Description").
				Placeholder("One year duration"),
			huh.NewInput().
				Value(&newConfig.Es.Validity.DurationString).
				Title("Validity Duration").
				Placeholder("365").
				Description("Duration in days"),
		),
	).WithAccessible(accessible)

	RunForm(EsForm)

	_ = spinner.New().Title("Sending Es connection config...").Accessible(accessible).Action(newConfig.SendEsConnection).Run()
}
