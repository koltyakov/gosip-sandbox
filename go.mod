module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.21
	github.com/Azure/go-autorest/autorest/adal v0.9.16
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.3 // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/howeyc/gopass v0.0.0-20210920133722-c8aef6fb66ef
	github.com/koltyakov/gosip v0.0.0-20211021174409-9642c3d37c35
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm v1.0.3
	github.com/vadimi/go-ntlm v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3 // indirect
	golang.org/x/sys v0.0.0-20211102061401-a2f17f7b995c // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)
