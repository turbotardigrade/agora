# Agora - Decentralized crowd-curated discussion platform

# Build

## Dependencies
*Note* that you need at least *go1.7* to run this project.

### Get IPFS from source and install it:
Steps summarized from [here](go get -u -d github.com/ipfs/go-ipfs).

```
go get -u -d github.com/ipfs/go-ipfs
cd $GOPATH/src/github.com/ipfs/go-ipfs
make install
```

### Get gx package manager and the rest of the dependencies
Gx is needed to get dependencies hosted on IPFS.

```
go get -u github.com/whyrusleeping/gx
go get ./...
```

## Build and Test

```
go build
go test
```

# Troubleshooting

## mdns lookup error

If you get this:

```
16:47:40.068 ERROR       mdns: mdns lookup error: failed to bind to any multicast udp port mdns.go:135
```

Fix it by raising the limit for the number of allowed file descriptors:

```
ulimit -n 5120
```