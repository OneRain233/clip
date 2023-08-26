# Clipboard Sync

## Description

This is a server for clipboard sync. It is written in Golang.

## Usage

### Server

```bash
go run main.go --mode [TCP|HTTP|BOTH]
```

### Client
If you use Gnome, you can use shortcut for get latest clipboard from server and send your clipboard to server.

```bash
go run client/client.go --api_host [Your Server Address] --mode http_get # Get latest clipboard from server
```

```bash
go run client/client.go --api_host [Your Server Address] --mode http_send # Send your clipboard to server
```

## TODO

- [ ] TCP Server