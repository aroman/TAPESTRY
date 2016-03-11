# tapestry

A set of tools to find public, user-uploaded videos taken of particular events, such as concerts, sporting events, and other occurrences with high spatial and temporal locality.

Currently, tapestry only supports mining and visualizing content from YouTube.


### mining for clusters

#### compiling the miner

1. `cd miner/`
2. install dependencies: `go get`
3. compile: `go build`

#### operating the miner

```
./mine --lat=51.5314303 --long=-0.128327 --radius=100km --before=02-05-2015 --terms="elton john london"
```

For a full reference of supported parameters, read the source (`main.go`).

### visualizing clusters

1. `cd miner/`
2. start the server: `go run serve.go database.go`
3. open a new terminal window
4. `cd frontend/`
5. install motion `npm -g install motion`
6. start the cluster explorer: `motion`
7. visit [http://localhost:4000](http://localhost:4000) in your browser

*coming soon*
