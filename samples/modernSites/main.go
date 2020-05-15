// Modern Sites Creation: https://docs.microsoft.com/en-us/sharepoint/dev/apis/site-creation-rest
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"

	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
)

// SiteCreationInfo ...
type SiteCreationInfo struct {
	Title               string `json:"Title"`
	URL                 string `json:"Url"`
	LCID                int    `json:"Lcid"`
	ShareByEmailEnabled bool   `json:"ShareByEmailEnabled"`
	Classification      string `json:"Classification"`
	Description         string `json:"Description"`
	WebTemplate         string `json:"WebTemplate"` // SITEPAGEPUBLISHING#0 - Communication Site, STS#3 - Team Site
	SiteDesignID        string `json:"SiteDesignId"`
	Owner               string `json:"Owner"`
}

// SiteCreationResponse ...
type SiteCreationResponse struct {
	SiteID     string `json:"SiteId"`
	SiteStatus int    `json:"SiteStatus"`
	// 0 - Not Found. The site doesn't exist.
	// 1 - Provisioning. The site is currently being provisioned.
	// 2 - Ready. The site has been created.
	// 3 - Error. An error occurred while provisioning the site.
	SiteURL string `json:"SiteUrl"`
}

func main() {
	strategy := flag.String("strategy", "saml", "Auth strategy")
	config := flag.String("config", "./config/private.json", "Config path")

	flag.Parse()

	// Binding auth & API client
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Index(strings.ToLower(authCnfg.GetSiteURL()), ".sharepoint.com") == -1 {
		log.Fatal(fmt.Errorf("SPSiteManager is only available in SPO"))
	}

	siteRandomPart := uuid.New().String()

	newSiteURL := fmt.Sprintf("%s/sites/commsite-%s", getHostURL(authCnfg.GetSiteURL()), siteRandomPart)
	newSiteName := "Communication Site " + siteRandomPart

	siteCreationInfo := &SiteCreationInfo{
		Title:        newSiteName,
		URL:          newSiteURL,
		LCID:         1033,
		WebTemplate:  "SITEPAGEPUBLISHING#0",                 // or "STS#3"
		SiteDesignID: "f6cc5403-0d63-442e-96c0-285923709ffc", // blank
	}

	// Creating new site
	creationResp, err := createSite(authCnfg, siteCreationInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Site creation status: %+v\n", creationResp)

	// Waiting for provisioning to complete
	for creationResp.SiteStatus != 2 {
		status, err := getSiteStatus(authCnfg, newSiteURL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Site status update: %+v\n", status)
		time.Sleep(time.Second * 2)
		creationResp.SiteStatus = status.SiteStatus
	}

	// Deleting the site
	if err := deleteSite(authCnfg, creationResp.SiteID); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Site is deleted: %s\n", newSiteURL)
}

func createSite(authCnfg gosip.AuthCnfg, siteCreationInfo *SiteCreationInfo) (*SiteCreationResponse, error) {
	sp := api.NewHTTPClient(&gosip.SPClient{AuthCnfg: authCnfg})

	endpoint := authCnfg.GetSiteURL() + "/_api/SPSiteManager/create"

	siteCreationBody := &struct {
		Request *SiteCreationInfo `json:"request"`
	}{
		Request: siteCreationInfo,
	}

	body, _ := json.Marshal(siteCreationBody)
	data, err := sp.Post(endpoint, bytes.NewBuffer(body), api.HeadersPresets.Minimalmetadata)
	if err != nil {
		return nil, err
	}

	resp := &SiteCreationResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func deleteSite(authCnfg gosip.AuthCnfg, siteID string) error {
	sp := api.NewHTTPClient(&gosip.SPClient{AuthCnfg: authCnfg})

	endpoint := authCnfg.GetSiteURL() + "/_api/SPSiteManager/delete"

	siteToDelete := &struct {
		SiteID string `json:"siteId"`
	}{}

	siteToDelete.SiteID = siteID

	body, _ := json.Marshal(siteToDelete)
	_, err := sp.Post(endpoint, bytes.NewBuffer(body), api.HeadersPresets.Minimalmetadata)
	return err
}

func getSiteStatus(authCnfg gosip.AuthCnfg, siteURL string) (*SiteCreationResponse, error) {
	sp := api.NewHTTPClient(&gosip.SPClient{AuthCnfg: authCnfg})

	endpoint := authCnfg.GetSiteURL() + "/_api/SPSiteManager/status?url='" + url.QueryEscape(siteURL) + "'"

	data, err := sp.Get(endpoint, api.HeadersPresets.Minimalmetadata)
	if err != nil {
		return nil, err
	}

	resp := &SiteCreationResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func getHostURL(absURL string) string {
	u := strings.Split(absURL, "/")
	return u[0] + "//" + u[2]
}
