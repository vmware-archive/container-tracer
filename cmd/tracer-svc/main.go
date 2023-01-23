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
	"strconv"
	"time"

	api "github.com/vmware-labs/container-tracer/api/svc"
	ctx "github.com/vmware-labs/container-tracer/internal/tracesvcctx"
)

var (
	description     = "container-tracer service"
	envAddress      = "TRACE_KUBE_API_ADDRESS"
	envVerbose      = "TRACE_KUBE_VERBOSE"
	envTracersPoll  = "TRACE_KUBE_DISCOVERY_POLL"
	envPodsSelector = "TRACE_KUBE_SELECTOR_PODS"
	envSvcSelector  = "TRACE_KUBE_SELECTOR_SVCS"

	defAddress     = ":8080"
	defPoll        = 10
	defPodSelector = "app=container-tracer-backend"
	defSvcSelector = "metadata.name=container-tracer-node"
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
	tracersPoll := flag.Int("poll", -1,
		fmt.Sprintf("Polling interval for tracers discovery, in seconds. Can be passed using %s environment variable as well.", envTracersPoll))
	cfg.Verbose = flag.Bool("verbose", false,
		fmt.Sprintf("Print informational logs on the standard output. Can be passed using %s environment variable as well.", envVerbose))
	cfg.PodSelector = flag.String("pods-selector", "",
		fmt.Sprintf("Label selector for filtering node tracer pods. Can be passed using %s environment variable as well.", envPodsSelector))
	cfg.SvcSelector = flag.String("svc-selector", "",
		fmt.Sprintf("Field selector for filtering node tracer services. Can be passed using %s environment variable as well.", envSvcSelector))

	flag.Parse()

	if *flApiAddr == "" {
		a := os.Getenv(envAddress)
		flApiAddr = &a
	}
	if *flApiAddr == "" {
		flApiAddr = &defAddress
	}

	if *cfg.PodSelector == "" {
		a := os.Getenv(envPodsSelector)
		cfg.PodSelector = &a
	}
	if *cfg.PodSelector == "" {
		cfg.PodSelector = &defPodSelector
	}

	if *cfg.SvcSelector == "" {
		a := os.Getenv(envSvcSelector)
		cfg.SvcSelector = &a
	}
	if *cfg.SvcSelector == "" {
		cfg.SvcSelector = &defSvcSelector
	}

	if *tracersPoll < 0 {
		if a, e := os.LookupEnv(envTracersPoll); e == true {
			if i, e := strconv.Atoi(a); e == nil {
				tracersPoll = &i
			} else {
				tracersPoll = &defPoll
			}
		} else {
			tracersPoll = &defPoll
		}
	}
	cfg.TracersPoll = time.Duration(*tracersPoll) * time.Second

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
		log.Fatal("Failed to create new container-tracer service: ", err)
		return
	}

	router := api.NewRouter(t)

	log.Printf("Listening for incoming API requests at %s", *addr)

	if err = router.Run(*addr); err != nil {
		log.Fatal("Failed to run the server:", err)
	}
}
