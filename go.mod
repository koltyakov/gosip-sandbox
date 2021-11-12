module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.22
	github.com/Azure/go-autorest/autorest/adal v0.9.17
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.9
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.4 // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/howeyc/gopass v0.0.0-20210920133722-c8aef6fb66ef
	github.com/koltyakov/gosip v0.0.0-20211112083007-3c824fb450b1
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm/v2 v2.4.1
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa // indirect
	golang.org/x/net v0.0.0-20211111160137-58aab5ef257a // indirect
	golang.org/x/sys v0.0.0-20211111213525-f221eed1c01e // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)
