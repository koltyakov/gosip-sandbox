// Package azurecert implements AAD Certificate Auth Flow
// See more:
//   - https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication
// Amongst supported platform versions are:
//   - SharePoint Online + Azure
package azurecert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - AAD Certificate Auth Flow
/* Config sample:
{
	"siteUrl": "https://contoso.sharepoint.com/sites/test",
	"tenantId": "e4d43069-8ecb-49c4-8178-5bec83c53e9d",
	"clientId": "628cc712-c9a4-48f0-a059-af64bdbb4be5",
	"certPath": "cert.pfx",
	"certPass": "password"
}
*/
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`  // SPSite or SPWeb URL, which is the context target for the API calls
	TenantID string `json:"tenantId"` // Azure Tenant ID
	ClientID string `json:"clientId"` // Azure Client ID
	CertPath string `json:"certPath"` // Azure certificate (.pfx) file location, relative to config location or absolute
	CertPass string `json:"certPass"` // Azure certificate export password

	authorizer  autorest.Authorizer
	privateFile string
	masterKey   string
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	c.privateFile = privateFile
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
	c.CertPath = path.Join(path.Dir(c.privateFile), c.CertPath)
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Decode(c.CertPass)
	if err == nil {
		c.CertPass = secret
	}
	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Encode(c.CertPass)
	if err != nil {
		return err
	}
	config := &AuthCnfg{
		SiteURL:  c.SiteURL,
		TenantID: c.TenantID,
		ClientID: c.ClientID,
		CertPath: c.CertPass,
		CertPass: secret,
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

	aadCnfg := auth.NewClientCertificateConfig(c.CertPath, c.CertPass, c.ClientID, c.TenantID)
	aadCnfg.Resource = resource
	authorizer, err := aadCnfg.Authorizer()
	if err != nil {
		return "", 0, err
	}

	c.authorizer = authorizer
	return "azure certificate via go-autorest/autorest/azure/auth", 0, nil
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "azurecert" }

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
