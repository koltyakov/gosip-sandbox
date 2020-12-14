module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.13
	github.com/Azure/go-autorest/autorest/adal v0.9.9
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.5
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/google/uuid v1.1.2
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/koltyakov/gosip v0.0.0-20201214145214-6c3f4747fe59
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm v1.0.3
	github.com/vadimi/go-ntlm v1.1.0 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
)
