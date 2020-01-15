# Dynamic authentication strategy

Currently, Gosip requires choosing a specific strategy or strategies to import and use within the application. However, some applications or utility tools might need more. There is no official universal resolver (yet it's planned), but in an application, some dynamic can be added on demand.

The sample shows a simple way of importing potentially demanded strategies and selecting one in runtime based on logic, CLI flags in the case of the sample.

## Implementation sample

Check out [code sources](./dynauth.go).

## Usage

```golang
package main

import (
	"flag"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
	"github.com/koltyakov/gosip/api"
)

var debug bool

func main() {

	strategy := flag.String("strategy", "ondemand", "Auth strategy")
	config := flag.String("config", "./config/private.json", "Config path")

	flag.Parse()

	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	// ...

}
```