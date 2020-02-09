# Folders & Files sync

The sample shows how to arrange local directory changes tracking with the corresponding synchronization with SharePoint document library.

The use-cases are:

- assets deployment pipelines
- file-based integrations

## Build

```bash
go build -o bin/filesync.exe ./samples/sync/
```

## Start process

Create `./config/private.json` with SAML auth credentials (or any other, but should be aligned with sources).

Run:

```bash
bin/filesync.exe -localFolder ./folder/to/watch -spFolder "Style Library"
```

Where:
- `-localFolder` is a local folder to watch
- `-spFolder` is SP folder to sync to

When applied changes are synced with SharePoint.

## All flags description

```bash
go run ./samples/sync/ -h
```

Flag | Description
-----|------------
`-strategy` | string, Auth strategy (default "saml")
`-config` | string, Config path (default "./config/private.json")
`-localFolder` | string, Local folder to watch
`-spFolder` | string, SP folder to sync to (default "SiteAssets")
`-skipSync` | bool, Skips initial sync of files on startup
`-watch` | bool, Watch local folder for changes
