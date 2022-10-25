// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	api "gitlab.eng.vmware.com/opensource/tracecruncher-api/api/svc"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracesvcctx"
)

var (
	description = "trace-kube service"
	envAddress  = "TRACE_KUBE_API_ADDRESS"
	envVerbose  = "TRACE_KUBE_VERBOSE"

	defAddress = ":8080"
)

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "%s: %s \n\n", os.Args[0], description)
	flag.PrintDefaults()
}

func getConfig() (*ctx.TraceKubeConfig, *string) {
	cfg := ctx.TraceKubeConfig{}

	flApiAddr := flag.String("address", "",
		fmt.Sprintf("IP address and port in format IP:port, used for listening for incoming API requests.Can be passed using %s environment variable as well",
			envAddress))
	cfg.Verbose = flag.Bool("verbose", false,
		fmt.Sprintf("Print informational logs on the standard output. Can be passed using %s environment variable as well.", envVerbose))

	flag.Parse()

	if *flApiAddr == "" {
		a := os.Getenv(envAddress)
		flApiAddr = &a
	}
	if *flApiAddr == "" {
		flApiAddr = &defAddress
	}
	if *cfg.Verbose == false {
		if _, ok := os.LookupEnv(envVerbose); ok {
			a := true
			cfg.Verbose = &a
		}
	}

	return &cfg, flApiAddr
}

func main() {
	var (
		t   *ctx.TraceKube
		err error
	)

	flag.Usage = usage

	cfg, addr := getConfig()
	if t, err = ctx.NewTraceKube(cfg); err != nil {
		log.Fatal("Failed to create new trace-kube service: ", err)
		return
	}

	router := api.NewRouter(t)

	log.Printf("Listening for incoming API requests at %s", *addr)

	if err = router.Run(*addr); err != nil {
		log.Fatal("Failed to run the server:", err)
	}
}
