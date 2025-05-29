package main

type SearchType struct {
	Entity string `json:"entity"`
	Type   string `json:"type"`
	State  bool   `json:"state"`
}

type GroupRight struct {
	Name             string   `json:"name"`
	Roles            []string `json:"roles"`
	AllowedUserDBIDs []int    `json:"allowed_userdb_ids"`
}
