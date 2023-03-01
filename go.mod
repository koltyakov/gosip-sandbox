module github.com/koltyakov/gosip-sandbox

go 1.16

// replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest v0.11.28
	github.com/Azure/go-autorest/autorest/adal v0.9.22
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.12
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/howeyc/gopass v0.0.0-20210920133722-c8aef6fb66ef
	github.com/koltyakov/gosip v0.0.0-20230225162407-d90f1db8d30e
	github.com/koltyakov/lorca v0.1.9-0.20230301032335-65834262f4bf
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm/v2 v2.4.1
	golang.org/x/net v0.7.0 // indirect
)
