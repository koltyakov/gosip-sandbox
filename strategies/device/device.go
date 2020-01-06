package device

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/koltyakov/gosip"
)

var tokenCache = map[string]*adal.ServicePrincipalToken{}

// AuthCnfg - AAD Device Flow auth config structure
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`  // SPSite or SPWeb URL, which is the context target for the API calls
	ClientID string `json:"clientId"` // Azure AD App Registration Client ID
	TenantID string `json:"tenantId"` // Azure AD App Registration Tenant ID
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	f, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	return json.Unmarshal(data, &c)
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	return nil
}

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, error) {
	u, _ := url.Parse(c.SiteURL)
	resource := fmt.Sprintf("https://%s", u.Host)

	// Check cached token per resource
	token := tokenCache[resource]
	if token != nil {
		// Return cached token if not expired
		if !token.Token().IsExpired() {
			return token.Token().AccessToken, nil
		}
		// Expired, try to refresh
		if err := token.Refresh(); err == nil {
			// Return refreshed token
			return token.Token().AccessToken, nil
		}
		// Failed to refresh, initiating for the device auth flow
	}

	// This is a sample, the prod Device auth flow would dump
	// access and refresh tokens to disk to avoid device wizard if possible

	config := auth.NewDeviceFlowConfig(c.ClientID, c.TenantID)
	config.Resource = resource

	token, err := config.ServicePrincipalToken()
	if err != nil {
		return "", err
	}

	tokenCache[resource] = token
	return token.Token().AccessToken, nil
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string {
	return c.SiteURL
}

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string {
	return "device"
}

// SetAuth authenticates request
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	accessToken, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return nil
}
