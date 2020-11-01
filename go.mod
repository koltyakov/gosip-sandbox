module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.10
	github.com/Azure/go-autorest/autorest/adal v0.9.5
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.3
	github.com/google/uuid v1.1.2
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/koltyakov/gosip v0.0.0-20201020124822-981a662a6e19
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm v1.0.3
	github.com/vadimi/go-ntlm v1.1.0 // indirect
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102 // indirect
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
)
