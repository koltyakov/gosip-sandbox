module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.13
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/google/uuid v1.2.0
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/koltyakov/gosip v0.0.0-20210522114515-0ffc4454bd9e
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm v1.0.3
	github.com/vadimi/go-ntlm v1.2.1 // indirect
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
)
