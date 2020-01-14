# Folders & Files sync

The sample shows how to arrange local directory changes tracking with corresponding syncronization with SharePoint document library.

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
bin/filesync.exe -watch ./folder/to/watch -spFolder "Style Library"
```

Where:
- `-watch` is a local folder to watch
- `-spFolder` is SP folder to sync to

When, applyed changes are synced with SharePoint.

## Sync all on start

Use `-syncAll` flag to upload local files on start.