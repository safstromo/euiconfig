package main

import (
	"fmt"
	"strconv"
	"time"
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
	daysInt, err := strconv.Atoi(v.DurationString)
	if err != nil {
		fmt.Println("Unable to convert duration string to int: ", v.DurationString)
	}

	duration := time.Duration(daysInt) * 24 * time.Hour
	v.Duration = duration.Milliseconds()
}
