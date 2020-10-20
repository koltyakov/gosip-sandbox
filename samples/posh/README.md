# Cookie consumption from within PowerShell

It's rare when you really need this but possible.

Let's assume a theoretical case, EdgeAccessCookie. When SharePoint is behind a WAP (Web Application Proxy) little or no drop-in libraries support this case, unfortunately. At least, it's happened to me always having trouble with WAP and libraries in the wild.

As a sort of a quick workaround is getting the Cookies and bypass them together with the requests.

## Build and pack

Go binaries usually weight around 12Mb and larger depending on used libraries. This is perfectly OK for most of the cases, however, for the embed purposes the less the bundle is better and there is reasoning to squeeze some juice out.

A great and simple approach is described in this [article](https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/).

### Build

`-ldflags="-s -w"` flags shrinks a binary almost twice a size. We're also disabling GC and the utility once started returns results and terminates, so disabling GC can improve performance a bit.

```bash
make build
```

### Pack

```bash
make pack
```

The `pack` task runs [upx](https://upx.github.io). Which archives a binary with appr. 30% ratio ending up with 2Mb binary which can be distributed to the consumers.

## Usage

### Invoke-WebRequest

[See more](./request.ps1)

### CSOM & PnP PoSH

[See more](./csom.ps1)