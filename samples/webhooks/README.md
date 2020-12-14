# Webhooks debug

## Exposing local dev server to cloud

- Install [ngrok](https://ngrok.com/download)
- `./ngrok authtoken <your_auth_token>`
- `./ngrok http 8080`

## Starting local dev server

```bash
./start.sh
```

## Reference

- [SharePoint list webhooks](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/lists/overview-sharepoint-list-webhooks)
- [Overview of SharePoint webhooks](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/overview-sharepoint-webhooks)
- [SharePoint webhooks sample reference implementation](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/webhooks-reference-implementation)
- [Using Azure Functions with SharePoint webhooks](https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/sharepoint-webhooks-using-azure-functions)