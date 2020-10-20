# Folders & Files upload

The sample shows how to arrange local directory upload to SharePoint document library.

The use-cases are:

- assets deployment pipelines
- documents migration

## Build

```bash
go build -o bin/spupload.exe ./samples/spupload/
```

## Start process

Create `./config/private.json` with SAML auth credentials (or any other, but should be aligned with sources).

Run:

```bash
bin/spsync.exe -localFolder ./upload/source/folder -spFolder "Shared Documents"
```

Where:
- `-localFolder` is a local folder to watch
- `-spFolder` is SP folder to sync to

When applied changes are synced with SharePoint.

## All flags description

```bash
go run ./samples/spupload/ -h
```

Flag | Description
-----|------------
`-strategy` | string, Auth strategy (default "saml")
`-config` | string, Config path (default "./config/private.json")
`-localFolder` | string, Local folder to watch
`-spFolder` | string, SP folder to sync to (default "Shared Documents")
`-concurrency` | int, a number of concurrent uploads (default 25)