package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"

	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
)

func main() {
	strategy := flag.String("strategy", "saml", "Auth strategy")
	config := flag.String("config", "./config/private.json", "Config path")

	flag.Parse()

	// Binding auth & API client
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatal(err)
	}
	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	// Getting all lists with role assigments details
	res, err := sp.Web().Lists().
		Select(`
			Title,
			HasUniqueRoleAssignments,
			RoleAssignments/Member/LoginName,
			RoleAssignments/Member/PrincipalType,
			RoleAssignments/RoleDefinitionBindings/Name,
			RoleAssignments/RoleDefinitionBindings/BasePermissions,
			RoleAssignments/RoleDefinitionBindings/RoleTypeKind
		`).
		Expand(`
			HasUniqueRoleAssignments,
			RoleAssignments/Member,
			RoleAssignments/RoleDefinitionBindings
		`).
		Get()

	if err != nil {
		log.Fatal(err)
	}

	var lists []*struct {
		Title                    string
		HasUniqueRoleAssignments bool
		RoleAssignments          []*api.RoleAssigment
	}

	// .Normalized() method aligns responses between different OData modes
	if err := json.Unmarshal(res.Normalized(), &lists); err != nil {
		log.Fatalf("unable to parse the response: %v", err)
	}

	for _, l := range lists {
		if !l.HasUniqueRoleAssignments {
			fmt.Printf("%s | inherits role assigments\n", l.Title)
		} else {
			fmt.Printf("%s | has unique role assigments\n", l.Title)
			for _, role := range l.RoleAssignments {
				fmt.Printf("  %s\n", role.Member.LoginName)
				for _, def := range role.RoleDefinitionBindings {
					fmt.Printf("    %s\n", def.Name)
				}
			}
		}
	}
}
