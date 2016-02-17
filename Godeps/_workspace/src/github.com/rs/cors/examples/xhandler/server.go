package main

import (
	"net/http"

	"github.com/aroman/tapestry-server/Godeps/_workspace/src/github.com/rs/cors"
	"github.com/aroman/tapestry-server/Godeps/_workspace/src/github.com/rs/xhandler"
	"github.com/aroman/tapestry-server/Godeps/_workspace/src/golang.org/x/net/context"
)

func main() {
	c := xhandler.Chain{}

	// Use default options
	c.UseC(cors.Default().HandlerC)

	mux := http.NewServeMux()
	mux.Handle("/", c.Handler(xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})))

	http.ListenAndServe(":8080", mux)
}
