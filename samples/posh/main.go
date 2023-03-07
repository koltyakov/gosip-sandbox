package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/koltyakov/gosip/auth/adfs"
)

func main() {
	siteURL := flag.String("siteUrl", "", "SharePoint Site Url")
	username := flag.String("username", "", "User login")
	password := flag.String("password", "", "User password")
	relyingParty := flag.String("relyingParty", "", "Relying party")
	adfsURL := flag.String("adfsUrl", "", "ADFS Url")
	adfsCookie := flag.String("adfsCookie", "", "ADFS Cookie")
	configPath := flag.String("configPath", "", "Connection config path")
	outFormat := flag.String("outFormat", "json", "Output Format: raw | json")

	flag.Parse()

	auth := &adfs.AuthCnfg{
		SiteURL:      *siteURL,
		Username:     *username,
		Password:     *password,
		RelyingParty: *relyingParty,
		AdfsURL:      *adfsURL,
		AdfsCookie:   *adfsCookie,
	}

	if *configPath != "" {
		err := auth.ReadConfig(*configPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "unable to read config: %v", err)
			os.Exit(1)
		}
	}

	authCookie, _, err := auth.GetAuth()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to authenticate: %v", err)
		os.Exit(1)
	}

	if authCookie == "" {
		_, _ = fmt.Fprint(os.Stderr, "can't get auth cookie")
		os.Exit(1)
	}

	if *outFormat == "raw" {
		_, _ = fmt.Fprint(os.Stdout, authCookie)
		os.Exit(0)
	}

	cookies := strings.Split(authCookie, "; ")

	json := "{"
	pcnt := 0
	for _, cookie := range cookies {
		c := strings.SplitN(cookie, "=", 2)
		if len(c) != 2 {
			continue
		}
		if pcnt > 0 {
			json += ","
		}
		json += fmt.Sprintf("\"%s\":\"%s\"", c[0], c[1])
		pcnt++
	}
	json += "}"

	_, _ = fmt.Fprint(os.Stdout, json)
}
