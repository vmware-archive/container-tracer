package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/opensource/tracecruncher-api/api"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracerctx"
)

var (
	apiAddr = flag.String("address", ":8080", "IP address and port in format IP:port, used for listening for incoming API requests.")
)

func main() {
	flag.Parse()

	var (
		t   *ctx.Tracer
		err error
	)
	if t, err = ctx.NewTracer(); err != nil {
		fmt.Println(err)
		return
	}

	router := api.NewRouter(t)

	log.Printf("Server started")

	router.Run(*apiAddr)
}
