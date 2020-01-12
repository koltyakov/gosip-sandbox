package ondemand

import (
	"strings"
	"time"

	"github.com/koltyakov/lorca"
)

// onDemandAuthFlow authenticates using On-Demand flow
func (c *AuthCnfg) onDemandAuthFlow(initialCookies *Cookies) (*Cookies, error) {
	startURL := "data:text/html,<html><head><title>Connecting to site: " + c.SiteURL + "</title><head></html>"
	ui, err := lorca.New(startURL, "", 480, 430)
	if err != nil {
		return nil, err
	}
	defer ui.Close()

	if initialCookies != nil {
		for _, cookie := range *initialCookies {
			ui.Send("Network.setCookie", cookie.toMap())
		}
	}

	ui.Load(c.SiteURL)

	currentURL := ""
	for strings.ToLower(currentURL) != strings.ToLower(c.SiteURL) {
		newURL := ui.Eval("window.location.href").String()
		if currentURL != newURL {
			// fmt.Printf("%s\n", newURL)
			currentURL = newURL
		}
		time.Sleep(500 * time.Microsecond)
	}
	resp := ui.Send("Network.getCookies", nil)
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	cookies := &Cookies{}
	if err := resp.Object()["cookies"].To(&cookies); err != nil {
		return nil, err
	}
	ui.Close()

	<-ui.Done()

	return cookies, nil
}
