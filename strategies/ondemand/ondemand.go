// Package ondemand implements On-Demand Auth Flow
// Amongst supported platform versions are:
//   - SharePoint Online
//   - SharePoint On-Premises (cookie-based auths)
package ondemand

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

var (
	cookieCache = map[string]*Cookies{}
	crypter     = cpass.Cpass("")
)

// AuthCnfg - On-Demand auth config structure
/* Config sample:
{
  "siteUrl": "https://contoso.sharepoint.com/sites/test",
}
*/
type AuthCnfg struct {
	SiteURL string `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
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
	config := &AuthCnfg{SiteURL: c.SiteURL}
	file, _ := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(privateFile, file, 0644)
}

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, error) {
	u, _ := url.Parse(c.SiteURL)

	// Check cached cookie per host
	cookies := cookieCache[u.Host]

	// Check disk cache
	if cookies == nil {
		cookies, _ = c.getCookieDiskCache()
	}

	if cookies != nil {
		// Return cached cookie if not expired
		if !cookies.isExpired() {
			return cookies.toString(), nil
		}
		// Expired, try to refresh
		cookies, err := c.onDemandAuthFlow(cookies)
		if err == nil {
			// Cache refreshed cookie
			c.cacheCookieToDisk(cookies)
			// Return refreshed token
			return cookies.toString(), nil
		}
		// Failed to refresh, initiating for the device auth flow
	}

	cookies, err := c.onDemandAuthFlow(nil)
	if err != nil {
		return "", err
	}

	c.cacheCookieToDisk(cookies)

	cookieCache[u.Host] = cookies
	return cookies.toString(), nil
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string {
	return c.SiteURL
}

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string {
	return "ondemand"
}

// SetAuth authenticates request
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	authCookie, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", authCookie)
	return nil
}

//=== File system cookie caching helpers ===//

// CleanCookieCache removes cookie information
func (c *AuthCnfg) CleanCookieCache() error {
	cookieCachePath := c.getCookieCachePath()
	u, err := url.Parse(c.SiteURL)
	if err != nil {
		return err
	}

	delete(cookieCache, u.Host)
	if err := os.Remove(cookieCachePath); err != nil {
		return err
	}
	return nil
}

// cacheCookieToDisk writes serialized cookies to temporary cache file
func (c *AuthCnfg) cacheCookieToDisk(cookies *Cookies) error {
	tmpDir := filepath.Join(os.TempDir(), "gosip")
	cookieCachePath := c.getCookieCachePath()

	cookieCache, err := json.Marshal(cookies)
	if err != nil {
		return err
	}
	cookieCacheE, _ := crypter.Encode(fmt.Sprintf("%s", cookieCache))
	cookieCache = []byte(cookieCacheE)

	os.MkdirAll(tmpDir, os.ModePerm)
	if err := ioutil.WriteFile(cookieCachePath, cookieCache, 0644); err != nil {
		return err
	}
	return nil
}

// getCookieDiskCache reads cookies from temporary cache file
func (c *AuthCnfg) getCookieDiskCache() (*Cookies, error) {
	cookieCachePath := c.getCookieCachePath()

	cookieCache, err := ioutil.ReadFile(cookieCachePath)
	if err != nil {
		return nil, err
	}
	cookieCacheD, _ := crypter.Decode(fmt.Sprintf("%s", cookieCache))
	cookieCache = []byte(cookieCacheD)

	cookies := &Cookies{}

	if err := json.Unmarshal(cookieCache, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

// getCookieCachePath gets local file system file path with token cache
func (c *AuthCnfg) getCookieCachePath() string {
	tmpDir := filepath.Join(os.TempDir(), "gosip")
	u, _ := url.Parse(c.SiteURL)
	return filepath.Join(tmpDir, c.GetStrategy()+"_"+u.Host)
}
