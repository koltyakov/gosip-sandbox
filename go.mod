module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.4
	github.com/Azure/go-autorest/autorest/adal v0.9.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.1
	github.com/google/uuid v1.1.2
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/koltyakov/gosip v0.0.0-20200903104039-7f86f0eaa40e
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm v1.0.3
	github.com/vadimi/go-ntlm v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/sys v0.0.0-20200905004654-be1d3432aa8f // indirect
)
