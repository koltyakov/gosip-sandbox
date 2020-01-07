# Cookie consumption from within PowerShell

It's rare when you really need this but possible.

Let's assume a theoretical case, EdgeAccessCookie. When SharePoint is behind a WAP (Web Application Proxy) little or no drop-in libraries support this case, unfortunately. At least, it's happened to me always having trouble with WAP and libraries in the wild.

As a sort of a quick workaround is getting the Cookies and bypass them together with the requests.