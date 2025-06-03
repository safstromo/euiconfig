package main

import (
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type WizepassFilter struct {
	Field       string `json:"field"`
	DisplayName string `json:"-"`
	State       bool   `json:"state"`
	Entity      string `json:"entity"`
}

var WIZEPASS_DTO_OPTIONS = []WizepassFilter{
	{
		Field:       "certificate_attributes",
		DisplayName: "Certificate attribute",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "certificate_serial_number",
		DisplayName: "Certificate serial number",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "created_on",
		DisplayName: "Created on",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "device_id",
		DisplayName: "Device ID",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "expires_on",
		DisplayName: "Expires on",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "notification_token",
		DisplayName: "Notification token",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "state",
		DisplayName: "State",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "user_id",
		DisplayName: "User ID",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "wizepass_id",
		DisplayName: "WizePass ID",
		State:       false,
		Entity:      "wizepass_dto",
	},
	{
		Field:       "tags",
		DisplayName: "Tags",
		State:       false,
		Entity:      "wizepass_dto",
	},
}

var WIZEPASS_ATTRIBUTE_OPTIONS = []WizepassFilter{
	{
		Field:       "common_name",
		DisplayName: "Common name",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "email_address",
		DisplayName: "Email address",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "given_name",
		DisplayName: "Given name",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "issuer",
		DisplayName: "Issuer",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "not_after",
		DisplayName: "Not after",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "not_before",
		DisplayName: "Not before",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "revoke_comment",
		DisplayName: "Revoke comment",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "serial_number",
		DisplayName: "Serial number",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "surname",
		DisplayName: "Surname",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "user_id",
		DisplayName: "User ID",
		State:       false,
		Entity:      "wizepass_attribute",
	},
	{
		Field:       "tags",
		DisplayName: "Tags",
		State:       false,
		Entity:      "wizepass_attribute",
	},
}

func FiltersForm(newConfig *Config) {
	Log.Info("Starting wizepass filter form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	// Reset to be able to run this again
	newConfig.SelectedAttributeFilters = []string{}
	newConfig.SelectedDTOFilters = []string{}

	dtoOptions := ConvertFiltersToHuhOptions(WIZEPASS_DTO_OPTIONS)
	attributeOptions := ConvertFiltersToHuhOptions(WIZEPASS_ATTRIBUTE_OPTIONS)

	filtersForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(dtoOptions...).
				Title("Wizepass DTO Filters").
				Value(&newConfig.SelectedDTOFilters),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(attributeOptions...).
				Title("Wizepass Attribute Filters").
				Value(&newConfig.SelectedAttributeFilters),
		),
	).WithAccessible(accessible)

	RunForm(filtersForm)

	_ = spinner.New().Title("Sending filters...").Accessible(accessible).Action(newConfig.SendFilters).Run()
}

func ConvertFiltersToHuhOptions(filters []WizepassFilter) []huh.Option[string] {
	Log.Info("Converting filters to options")
	options := make([]huh.Option[string], len(filters))
	for i, filter := range filters {
		option := huh.NewOption(filter.DisplayName, filter.Field)
		if filter.State {
			option = option.Selected(true)
		}
		options[i] = option
	}
	return options
}

func GetFilters(allFilters []WizepassFilter, selectedFields []string) []WizepassFilter {
	Log.Info("Converting selcted fields to filters")
	selectedMap := make(map[string]bool)
	for _, field := range selectedFields {
		selectedMap[field] = true
	}

	var activeFilters []WizepassFilter
	for _, filter := range allFilters {
		if selectedMap[filter.Field] {
			f := filter
			f.State = true
			activeFilters = append(activeFilters, f)
		}
	}
	return activeFilters
}
