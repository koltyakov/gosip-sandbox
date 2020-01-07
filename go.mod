module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest/adal v0.8.1
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/koltyakov/gosip v0.0.0-20200107164051-f4efd3d97179
	golang.org/x/sys v0.0.0-20200107162124-548cf772de50 // indirect
)
