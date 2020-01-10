# Alternative NTLM client

See [more](https://github.com/koltyakov/gosip/issues/14).

## Auth configuration and usage

```golang
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip-sandbox/strategies/ntlm"
)

func main() {

	authCnfg := &strategy.AuthCnfg{
		SiteURL:  os.Getenv("SPAUTH_SITEURL"),
		Username: os.Getenv("SPAUTH_USERNAME"),
    Domain:   os.Getenv("SPAUTH_DOMAIN"),
		Password: os.Getenv("SPAUTH_Password"),
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	res, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Site title: %s\n", res.Data().Title)

}
```