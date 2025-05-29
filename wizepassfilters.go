package main

type WizepassFilter struct {
	Field       string `json:"type"`
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
