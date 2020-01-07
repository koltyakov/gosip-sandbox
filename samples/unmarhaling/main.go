package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"
)

func main() {
	// Binding auth & API client
	configPath := "./config/private.saml.json"
	authCnfg := &strategy.AuthCnfg{}
	if err := authCnfg.ReadConfig(configPath); err != nil {
		log.Fatalf("unable to get config: %v", err)
	}
	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	// Getting items from a custom list
	list := sp.Web().GetList("Lists/MyList")
	data, err := list.Items().Select("Id,CustomField").Get()
	if err != nil {
		log.Fatalln(err)
	}

	// Define a stuct or map[string]interface{} for unmarshalling
	items := []*struct {
		ID     int    `json:"Id"`
		Custom string `json:"CustomField"`
	}{}

	// Use `data.Normalized()` or `data, _ = api.NormalizeODataCollection(data)`
	// for OData modes normalization.

	// For a single item, use
	// `data.Normalized()` or `data = api.NormalizeODataItem(data)

	// .Normalized() method aligns responses between different OData modes
	if err := json.Unmarshal(data.Normalized(), &items); err != nil {
		log.Fatalf("unable to parse the response: %v", err)
	}

	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}

}
