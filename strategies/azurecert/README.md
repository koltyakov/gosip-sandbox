# feat(azurecert): Auto-detect Azure AD endpoint from SharePoint URL

## Summary

Auto-detects the correct Azure AD endpoint for sovereign clouds based on the SharePoint site URL domain suffix. No configuration changes required.

## How it works

The `GetAuth()` function now automatically selects the correct Azure AD endpoint:

| SharePoint Domain | Azure AD Endpoint |
|-------------------|-------------------|
| `*.sharepoint.com` | `login.microsoftonline.com` (Commercial) |
| `*.sharepoint.us` | `login.microsoftonline.us` (GCC High) |
| `*.sharepoint.cn` | `login.chinacloudapi.cn` (China) |
| `*.sharepoint.de` | `login.microsoftonline.de` (Germany) |

## Usage

No changes to configuration required. Just use the SharePoint site URL:

**Commercial:**

```json
{
  "siteUrl": "https://contoso.sharepoint.com/sites/test",
  "tenantId": "...",
  "clientId": "...",
  "certPath": "cert.pfx",
  "certPass": "password"
}
```

**GCC High:**

```json
{
  "siteUrl": "https://contoso.sharepoint.us/sites/test",
  "tenantId": "...",
  "clientId": "...",
  "certPath": "cert.pfx",
  "certPass": "password"
}
```

## Backwards Compatible

- Commercial cloud (`.sharepoint.com`) continues to work exactly as before
- No breaking changes to existing configurations
- Uses Azure SDK's built-in cloud environment constants
