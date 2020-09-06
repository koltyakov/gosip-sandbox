# Azure AD Environment-Based Auth Flow Sample

The sample shows Gosip [custom auth](https://go.spflow.com/auth/custom-auth) with [AAD Environment-Based Authorization](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication).

## Custom auth implementation

Checkout [the code](./azureenv.go).

## Azure App registration

1\. Create or use existing app registration

2\. Make sure that the app is configured for a specific auth scenario:
- Client credentials
- Certificate
- Username/Password
- Managed identity

## Auth configuration and usage

```golang
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip-sandbox/strategies/azureenv"
)

func main() {

	authCnfg := &strategy.AuthCnfg{}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	res, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Site title: %s\n", res.Data().Title)

}
```
