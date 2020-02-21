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
	strategy := flag.String("strategy", "", "Auth strategy")
	config := flag.String("config", "", "Config path")

	flag.Parse()

	// Binding auth & API client
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatal(err)
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
	var items []*struct {
		ID     int    `json:"Id"`
		Custom string `json:"CustomField"`
	}

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
