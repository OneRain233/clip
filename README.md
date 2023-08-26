# Clipboard Sync

## Description

This is a server for clipboard sync. It is written in Golang.

## Usage

### Server

```bash
go run main.go --mode [TCP|HTTP|BOTH]
```

### Client

#### Linux

If you use Gnome, you can use shortcut for get latest clipboard from server and send your clipboard to server.

```bash
go run client/client.go --api_host [Your Server Address] --mode http_get # Get latest clipboard from server
```

```bash
go run client/client.go --api_host [Your Server Address] --mode http_send # Send your clipboard to server
```


#### Mobile

##### Android

You can use Tasker to send your clipboard to server via API.
`http://[Your Server Address]:[Your Server Port]/clipboard/add`

And you can use Tasker to get latest clipboard from server via API.
`http://[Your Server Address]:[Your Server Port]/clipboard/latest`

##### iOS
You can use Automator to send your clipboard to server via API like Android.

## TODO

- [ ] TCP Server
- [ ] Different OS support
- [ ] Client for different OS