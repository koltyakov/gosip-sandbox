# Local Dev Proxy

## Use-cases

During front-end development, it is common starting a local dev server with the application and API server. The sample shows an essential local proxy for bypassing anonymous requests from a workbench to a SharePoint server. That can save lots of time spent on mocking responses up.

Another scenario could be to put such a server between SharePoint and 3rd party application, which for some reason doesn't support some authentication mechanism, and it's not possible to inject this auth into a tool.

## Start

```bash
go run ./samples/proxy -strategy adfs -config ./config/private.json -port 9090
```

## HTTPS

```bash
openssl genrsa -out ./config/certs/private.key 2048
openssl req -new -x509 -sha256 -key ./config/certs/private.key -out ./config/certs/public.crt -days 3650
```

```bash
go run ./samples/proxy -strategy adfs -config ./config/private.onprem-wap-adfs.json -port 443 -sslKey ./config/certs/private.key -sslCert ./config/certs/public.crt
```