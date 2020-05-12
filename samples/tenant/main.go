package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
	"github.com/koltyakov/gosip/api"
)

func main() {

	strategy := flag.String("strategy", "saml", "Auth strategy")
	config := flag.String("config", "./config/private.json", "Config path")

	flag.Parse()

	// Binding auth & API client
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatalln(err)
	}

	if strings.Index(strings.ToLower(authCnfg.GetSiteURL()), ".sharepoint.com") == -1 {
		log.Fatalln(fmt.Errorf("is only compatible with SPO"))
	}

	// spClient := api.NewHTTPClient(&gosip.SPClient{AuthCnfg: authCnfg})
	web := api.NewWeb(&gosip.SPClient{AuthCnfg: authCnfg}, getTenantURL(authCnfg.GetSiteURL())+"/_api/web", nil)

	lists, err := web.Lists().Select("Id,Title").OrderBy("Title", true).Get()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("CA site lists:")
	for _, list := range lists.Data() {
		d := list.Data()
		fmt.Printf(" - %s\n", d.Title)
	}

	allSiteCollections, err := web.Lists().GetByTitle("DO_NOT_DELETE_SPLIST_TENANTADMIN_AGGREGATED_SITECOLLECTIONS").Items().Top(5000).Get()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("All site collections:")
	for _, sc := range allSiteCollections.Data() {
		d := sc.Data()
		fmt.Printf(" - %s\n", d.Title)
	}
}

func getTenantURL(absURL string) string {
	u := strings.Split(absURL, "/")
	tenantURL := u[0] + "//" + u[2]
	tenantURL = strings.Replace(tenantURL, "-admin.", ".", -1)
	tenantURL = strings.Replace(tenantURL, ".sharepoint.com", "-admin.sharepoint.com", -1)
	return tenantURL
}
