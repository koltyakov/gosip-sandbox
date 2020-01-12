module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	github.com/Azure/go-autorest/autorest/adal v0.8.1
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/ThomsonReutersEikon/go-ntlm v0.0.0-20181130171125-cf23bd1ecf18 // indirect
	github.com/koltyakov/gosip v0.0.0-20200109060436-2111b2e6d6bd
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/vadimi/go-http-ntlm v1.0.1
	golang.org/x/crypto v0.0.0-20200109152110-61a87790db17 // indirect
	golang.org/x/sys v0.0.0-20200107162124-548cf772de50 // indirect
)
