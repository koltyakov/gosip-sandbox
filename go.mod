module github.com/koltyakov/gosip-sandbox

go 1.13

replace github.com/koltyakov/gosip-sandbox => ./

require (
	fyne.io/fyne v1.3.0
	github.com/Azure/go-autorest/autorest v0.10.1 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/ThomsonReutersEikon/go-ntlm v0.0.0-20181130171125-cf23bd1ecf18 // indirect
	github.com/google/uuid v1.1.1
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/koltyakov/gosip v0.0.0-20200517122539-b234af4b5d8c
	github.com/koltyakov/lorca v0.1.9-0.20200112132759-701f901adf53
	github.com/radovskyb/watcher v1.0.7
	github.com/vadimi/go-http-ntlm v1.0.1
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
)
