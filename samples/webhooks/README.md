# SharePoint Webhooks with Golang

## Exposing local dev server to cloud

- Install [ngrok](https://ngrok.com/download)
- `./ngrok authtoken <your_auth_token>`
- `./ngrok http 8080`

The resource for Weebhooks must be exposed via HTTPS on a public IP. ngrok tunnels locally started service throught its infrasctructure.

## Environment variables

Create `.env` file with, or set up the following environment variables:

Variable | Description
---------|------------
SPAUTH_SITEURL      | SharePoint Site URL
SPAUTH_CLIENTID     | Add-In Client ID
SPAUTH_CLIENTSECRET | Add-In Client Secret
NOTIFICATIONS_URL   | `https://{host}/api/notifications`, optional for local debug, then deployed, e.g. as [Azure Function](https://github.com/koltyakov/az-fun-go-sp) should be provided explicitly.

See [more details](https://go.spflow.com/auth/strategies/addin). Add-In Only Auth strategy is one of the possible options.

## Starting local dev server

```bash
./ngrok http 8080
./start.sh
```

The `start.sh` script detects ngrok public endpoint automatically, compiles and lift the dev server, ready to receive webhooks requests.

## Subscribing webhooks to a list

```bash
curl http://localhost:8080/api/subscribe?listName=Site%20Pages
```

`listName` parameter receives the display name of an existing list.

The hook is registered for a short period of time (10 mins, within the sample). Production grade webhook expiration should be in 6 months margin.

Apply any changes to the list items, after a while the service should receive some updates. Updates are printed in the console.

## Reference

- [SharePoint list webhooks](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/lists/overview-sharepoint-list-webhooks)
- [Overview of SharePoint webhooks](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/overview-sharepoint-webhooks)
- [SharePoint webhooks sample reference implementation](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/webhooks-reference-implementation)
- [Using Azure Functions with SharePoint webhooks](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/sharepoint-webhooks-using-azure-functions)