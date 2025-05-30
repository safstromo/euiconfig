package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type Userdb struct {
	Id   int    `json:"-"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func UserDbConnectionForm(newConfig *Config) {
	Log.Info("Starting userdb connection form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	userDb := Userdb{}

	userdbForm := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Userdb connection"),

			huh.NewInput().
				Value(&userDb.Name).
				Title("Please enter Name").
				Placeholder("Defaut Userdb"),
			huh.NewInput().
				Value(&userDb.Url).
				Title("Please enter Userdb Url").
				Placeholder("http://localhost:8081").
				Description("The url to the Userdb to configure"),
		),
	).WithAccessible(accessible)

	err := userdbForm.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	newConfig.AddedUserDbs = append(newConfig.AddedUserDbs, userDb)

	_ = spinner.New().Title("Sending Userdb connection config...").Accessible(accessible).Action(newConfig.SendUserdbConnection).Run()
}

type UserDBConfig struct {
	UseLdap                      bool     `json:"use_ldap"`
	LdapTLS                      bool     `json:"ldap_tls"`
	LdapURL                      string   `json:"ldap_url"`
	LdapPort                     int      `json:"ldap_port"`
	LdapBaseDn                   string   `json:"ldap_base_dn"`
	LdapUseAlternativeAttributes bool     `json:"ldap_use_alternative_attributes"`
	LdapAdditionalSearchFilters  []string `json:"ldap_additional_searchfilters"`
	SearchResultLimit            int      `json:"search_result_limit"`
	MapLdapUserID                string   `json:"map_ldap_user_id"`
	MapLdapCommonName            string   `json:"map_ldap_common_name"`
	MapLdapGivenName             string   `json:"map_ldap_given_name"`
	MapLdapSurname               string   `json:"map_ldap_surname"`
	MapLdapEmail                 string   `json:"map_ldap_email"`
	MapLdapPersonalNumber        string   `json:"map_ldap_personal_number"`
	MapLdapUserPrincipalName     string   `json:"map_ldap_user_principal_name"`
}

func NewDefaultUserDBConfig() UserDBConfig {
	return UserDBConfig{
		UseLdap:                      true,
		LdapTLS:                      false,
		LdapURL:                      "localhost",
		LdapPort:                     389,
		LdapBaseDn:                   "dc=wizepass,dc=com",
		LdapUseAlternativeAttributes: false,
		LdapAdditionalSearchFilters:  []string{"objectclass=inetOrgPerson"},
		SearchResultLimit:            200,
		MapLdapUserID:                "userPrincipalName",
		MapLdapCommonName:            "cn",
		MapLdapGivenName:             "givenName",
		MapLdapSurname:               "sn",
		MapLdapEmail:                 "mail",
		MapLdapPersonalNumber:        "personalNumber",
		MapLdapUserPrincipalName:     "userPrincipalName",
	}
}

func UserDbConfigForm(newConfig *Config) {
	Log.Info("Starting userdb config form")
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	searchfilters := "objectclass=person"
	searchlimit := "200"
	ldapPort := "389"

	userdbForm := huh.NewForm(

		huh.NewGroup(huh.NewNote().
			Title("Userdb config"),

			huh.NewInput().
				Value(&newConfig.UserDBConfig.LdapURL).
				Title("Please enter LDAP URL/IP").
				Placeholder("http://localhost.com"),
			huh.NewInput().
				Value(&ldapPort).
				Title("Please enter LDAP port").
				Placeholder("389").
				Validate(func(str string) error {
					_, err := strconv.Atoi(str)
					if err != nil {
						return errors.New("Unable to convert port string")
					}
					return nil
				}),
			huh.NewInput().
				Value(&searchlimit).
				Title("Search Result limin").
				Placeholder("200").
				Description("Limits the number of results in a search").
				Validate(func(str string) error {
					_, err := strconv.Atoi(str)
					if err != nil {
						return errors.New("Unable to convert search limit string")
					}
					return nil
				}),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.LdapBaseDn).
				Title("LDAP base DN").
				Placeholder("dc=lab2019,dc=local"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapGivenName).
				Title("Map LDAP given name").
				Placeholder("givenName"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapCommonName).
				Title("Map LDAP common name").
				Placeholder("cn"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapSurname).
				Title("Map LDAP surname").
				Placeholder("sn"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapUserID).
				Title("Map LDAP user id").
				Placeholder("samaccountname"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapEmail).
				Title("Map LDAP email").
				Placeholder("mail"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapPersonalNumber).
				Title("Map LDAP personal number").
				Placeholder("personalNumber"),
			huh.NewInput().
				Value(&newConfig.UserDBConfig.MapLdapUserPrincipalName).
				Title("Map LDAP principal name").
				Placeholder("userPrincipalName"),
			huh.NewInput().
				Value(&searchfilters).
				Title("LDAP additional seach filters").
				Placeholder("objectclass=person,objectclass=inetOrgPerson").
				Description("comma separated list"),
			huh.NewConfirm().
				Title("Use LDAP").
				Value(&newConfig.UserDBConfig.UseLdap).
				Affirmative("Yes!").
				Negative("No.").
				Description("Use ldap or database"),
			huh.NewConfirm().
				Title("Use LDAP TLS").
				Value(&newConfig.UserDBConfig.LdapTLS).
				Affirmative("Yes!").
				Negative("No.").
				Description("Use LDAPS connection"),
			huh.NewConfirm().
				Title("LDAP use alternative attributes").
				Value(&newConfig.UserDBConfig.LdapUseAlternativeAttributes).
				Affirmative("Yes!").
				Negative("No.").
				Description("Enable use of alternative attibutes"),
		),
	).WithAccessible(accessible)

	err := userdbForm.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	newConfig.UserDBConfig.LdapAdditionalSearchFilters = strings.Split(searchfilters, ",")

	port, _ := strconv.Atoi(ldapPort)
	newConfig.UserDBConfig.LdapPort = port

	limit, _ := strconv.Atoi(searchlimit)
	newConfig.UserDBConfig.SearchResultLimit = limit

	_ = spinner.New().Title("Sending Userdb connection config...").Accessible(accessible).Action(newConfig.SendUserDbConfig).Run()
}
