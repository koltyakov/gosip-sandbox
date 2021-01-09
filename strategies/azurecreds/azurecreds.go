// Package azurecreds implements AAD Username/Password Auth Flow
// See more:
//   - https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication
// Amongst supported platform versions are:
//   - SharePoint Online + Azure
package azurecreds

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
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - AAD Username/Password Auth Flow
// To use this strategy public client flows mobile and desktop should be enabled in the app registration
/* Config sample:
{
	"siteUrl": "https://contoso.sharepoint.com/sites/test",
	"tenantId": "e4d43069-8ecb-49c4-8178-5bec83c53e9d",
	"clientId": "628cc712-c9a4-48f0-a059-af64bdbb4be5",
	"username": "user@contoso.com",
	"password": "password"
}
*/
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`  // SPSite or SPWeb URL, which is the context target for the API calls
	TenantID string `json:"tenantId"` // Azure Tenant ID
	ClientID string `json:"clientId"` // Azure Client ID
	Username string `json:"username"` // AAD user name
	Password string `json:"password"` // AAD user password

	authorizer autorest.Authorizer
	masterKey  string
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	f, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	byteValue, _ := ioutil.ReadAll(f)
	return c.ParseConfig(byteValue)
}

// ParseConfig parses credentials from a provided JSON byte array content
func (c *AuthCnfg) ParseConfig(byteValue []byte) error {
	if err := json.Unmarshal(byteValue, &c); err != nil {
		return err
	}
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Decode(c.Password)
	if err == nil {
		c.Password = secret
	}
	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Encode(c.Password)
	if err != nil {
		return err
	}
	config := &AuthCnfg{
		SiteURL:  c.SiteURL,
		TenantID: c.TenantID,
		ClientID: c.ClientID,
		Username: c.Username,
		Password: secret,
	}
	file, _ := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(privateFile, file, 0644)
}

// SetMasterkey defines custom masterkey
func (c *AuthCnfg) SetMasterkey(masterKey string) { c.masterKey = masterKey }

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, int64, error) {
	u, _ := url.Parse(c.SiteURL)
	resource := fmt.Sprintf("https://%s", u.Host)

	aadCnfg := auth.NewUsernamePasswordConfig(c.Username, c.Password, c.ClientID, c.TenantID)
	aadCnfg.Resource = resource
	authorizer, err := aadCnfg.Authorizer()
	if err != nil {
		return "", 0, err
	}

	c.authorizer = authorizer
	return "azure username/password via go-autorest/autorest/azure/auth", 0, nil
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "azurecreds" }

// SetAuth authenticates request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	if _, _, err := c.GetAuth(); err != nil {
		return err
	}
	_, err := c.authorizer.WithAuthorization()(preparer{}).Prepare(req)
	return err
}

// Preparer implements autorest.Preparer interface
type preparer struct{}

// Prepare satisfies autorest.Preparer interface
func (p preparer) Prepare(req *http.Request) (*http.Request, error) { return req, nil }
