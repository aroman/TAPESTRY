# TAPESTRY

A set of tools to find public, user-uploaded videos taken of particular events, such as concerts, sporting events, and other occurrences with high spatial and temporal locality.

Currently, tapestry only supports mining from YouTube.

### mining for clusters

#### compiling the YouTube miner

1. install dependencies: `go get`
2. `go build -o mine-youtube miners/youtube.go`
3. `./mine-youtube`

#### operating the miner

To run the miner one-off, you can invoke it like this:
```
$ ./mine --lat=51.5314303 --long=-0.128327 --radius=100km --before=02-05-2015 --terms="elton john london"
```

You probably don't want to run it by hand, though.

Connect the event generator to the miner like so:
```
$ cd event-generator
$ python generate.py | uniq | tr '\n' '\0' | xargs -0 -n1  ../mine-youtube
```

For a full reference of supported parameters, read the source (`miners/youtube.go`).

### exploring clusters

1. install dependencies: `go get`
2. start the server: `go run serve.go`
3. visit [http://localhost:8000](http://localhost:8000) in your browser
