package main

type GroupRight struct {
	Name             string   `json:"name"`
	Roles            []string `json:"roles"`
	AllowedUserDBIDs []int    `json:"allowed_userdb_ids"`
}
