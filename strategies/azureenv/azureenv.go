// Package azureenv implements AAD Environment-Based Auth Flow
// See more: https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication
// Amongst supported platform versions are:
//   - SharePoint Online + Azure
// Azure Environment-Based supported strategies:
//   - Client credentials
//   - Certificate
//   - Username/Password
//   - Managed identity
package azureenv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/koltyakov/gosip"
)

// AuthCnfg - AAD Environment-Based Auth Flow
// No creds settings in auth config are required, the configuration is happening through environment variables:
// https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication
/* Config sample:
{ "siteUrl": "https://contoso.sharepoint.com/sites/test" }
*/
type AuthCnfg struct {
	SiteURL    string `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
	authorizer autorest.Authorizer
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	f, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	data, _ := ioutil.ReadAll(f)
	return json.Unmarshal(data, &c)
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	config := &AuthCnfg{SiteURL: c.SiteURL}
	file, _ := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(privateFile, file, 0644)
}

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, error) {
	u, _ := url.Parse(c.SiteURL)
	resource := fmt.Sprintf("https://%s", u.Host)

	authorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(resource)
	if err != nil {
		return "", err
	}

	c.authorizer = authorizer

	return "azure environment via go-autorest/autorest/azure/auth", nil
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string {
	return c.SiteURL
}

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string {
	return "azureenv"
}

// SetAuth authenticates request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	if _, err := c.GetAuth(); err != nil {
		return err
	}
	req, err := c.authorizer.WithAuthorization()(preparer{}).Prepare(req)
	return err
}

// Preparer implements autorest.Preparer interface
type preparer struct{}

// Prepare satisfies autorest.Preparer interface
func (p preparer) Prepare(req *http.Request) (*http.Request, error) {
	return req, nil
}
