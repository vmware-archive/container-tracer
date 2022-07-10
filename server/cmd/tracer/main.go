package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/opensource/tracecruncher-api/api"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracer"
)

var (
	apiAddr = flag.String("addres", ":8080", "IP address and port in format IP:port, used for listening for incoming API requests.")
)

func main() {
	flag.Parse()

	var (
		t   *tracer.Tracer
		err error
	)
	if t, err = tracer.NewTracer(); err != nil {
		fmt.Println(err)
		return
	}

	router := api.NewRouter(t)

	log.Printf("Server started")

	router.Run(*apiAddr)
}
