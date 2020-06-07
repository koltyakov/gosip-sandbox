package main

import (
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
	list := sp.Web().GetList("Lists/Records")
	err = list.Items().GetByID(1).Records().Declare()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done")

}
