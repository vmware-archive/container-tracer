// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Enyinna Ochulor <eochulor@vmware.com>
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 */

/*
   -address string
		IP address and port in format IP:port, used for listening for incoming API requests.
		Can be passed using TRACER_API_ADDRESS environment variable as well
  -cri-endpoint string
		Path to the CRI endpoint. Can be passed using TRACER_CRI_ENDPOINT environment variable as well.
  -jaeger-endpoint string
		URL or name of the jaeger endpoint service, used to send collected traces.
		Can be passed using TRACER_JEAGER_ENDPOINT environment variable as well.
  -node-name string
		Name of the node, which runs that tracer instance.
		Can be passed using TRACER_NODE_NAME environment variable as well.
  -pod-name string
		Name of the tracer pod, used to verify the CRI endpoint.
		Can be passed using TRACER_POD_NAME environment variable as well.
  -procfs-path string
		Path to the /proc fs mount point. Can be passed using TRACER_PROCFS_PATH environment variable as well.
  -run-path value
		Path to the run directories, to look for cri endpoints.
		Can be passed using TRACER_RUN_PATHS environment variable as well.
  -sysfs-path string
		Path to the /sys fs mount point. Can be passed using TRACER_SYSFS_PATH environment variable as well.
  -trace-hooks string
		Location of the directory with trace helper applications.
		Can be passed using TRACER_HOOKS environment variable as well.
  -use-procfs
		Force using procfs for containers discovery, even if CRI is available.
		Can be passed using TRACER_FORCE_PROCFS environment variable as well.
  -verbose
		Print informational logs on the standard output.
		Can be passed using TRACER_VERBOSE environment variable as well.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/vmware-labs/container-tracer/api/node"
	"github.com/vmware-labs/container-tracer/internal/logger"
	"github.com/vmware-labs/container-tracer/internal/pods"
	hooks "github.com/vmware-labs/container-tracer/internal/tracehook"
	trace "github.com/vmware-labs/container-tracer/internal/tracerctx"
)

var (
	appName     = "container-tracer"
	description = "Trace containers running on the local node."
	envAddress  = "TRACER_API_ADDRESS"
	envVerbose  = "TRACER_VERBOSE"
	envNodeName = "TRACER_NODE_NAME"

	defAddress = ":8080"
)

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "%s: %s \n\n", os.Args[0], description)
	flag.PrintDefaults()
}

type stringsFlag []string

func (arr *stringsFlag) String() string {
	str := ""
	if arr != nil {
		for _, s := range *arr {
			str += s + " "
		}
	}
	return str
}

func (arr *stringsFlag) Set(s string) error {

	sep := func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t'
	}
	all := strings.FieldsFunc(s, sep)
	for _, a := range all {
		sa := strings.Trim(a, " \t\n")
		if sa != "" {
			*arr = append(*arr, sa)
		}
	}

	return nil
}

func getConfig() (*trace.TracerConfig, *string) {
	var runPathsArg stringsFlag
	cfg := trace.TracerConfig{}

	flApiAddr := flag.String("address", "",
		fmt.Sprintf("IP address and port in format IP:port, used for listening for incoming API requests.Can be passed using %s environment variable as well",
			envAddress))
	cfg.Verbose = flag.Bool("verbose", false,
		fmt.Sprintf("Print informational logs on the standard output. Can be passed using %s environment variable as well.", envVerbose))
	cfg.NodeName = flag.String("node-name", "",
		fmt.Sprintf("Name of the node, which runs that tracer instance. Can be passed using %s environment variable as well.", envNodeName))

	cfg.Hook.Procfs = flag.String("procfs-path", "",
		fmt.Sprintf("Path to the /proc fs mount point. Can be passed using %s environment variable as well.", hooks.EnvProcfs))
	cfg.Hook.Sysfs = flag.String("sysfs-path", "",
		fmt.Sprintf("Path to the /sys fs mount point. Can be passed using %s environment variable as well.", hooks.EnvSysfs))
	cfg.Hook.HooksPath = flag.String("trace-hooks", "",
		fmt.Sprintf("Location of the directory with trace helper applications. Can be passed using %s environment variable as well.", hooks.EnvHooks))

	flag.Var(&runPathsArg, "run-path",
		fmt.Sprintf("Path to the run directories, to look for cri endpoints. Can be passed using %s environment variable as well.", pods.EnvRunPaths))
	cfg.Pod.Cri.Endpoint = flag.String("cri-endpoint", "",
		fmt.Sprintf("Path to the CRI endpoint. Can be passed using %s environment variable as well.", pods.EnvCri))
	cfg.Pod.Cri.PodName = flag.String("pod-name", "",
		fmt.Sprintf("Name of the tracer pod, used to verify the CRI endpoint. Can be passed using %s environment variable as well.", pods.EnvPodName))
	cfg.Pod.ForceProc = flag.Bool("use-procfs", false,
		fmt.Sprintf("Force using procfs for containers discovery, even if CRI is available. Can be passed using %s environment variable as well.", pods.EnvForceProcfs))

	cfg.Logger.JaegerEndpoint = flag.String("jaeger-endpoint", "",
		fmt.Sprintf("URL or name of the jaeger endpoint service, used to send collected traces. Can be passed using %s environment variable as well.", logger.EnvLoggerJaegerEndpoint))

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
	if *cfg.NodeName == "" {
		a := os.Getenv(envNodeName)
		cfg.NodeName = &a
	}

	if *cfg.Hook.Procfs == "" {
		a := os.Getenv(hooks.EnvProcfs)
		cfg.Hook.Procfs = &a
	}
	if *cfg.Hook.Sysfs == "" {
		a := os.Getenv(hooks.EnvSysfs)
		cfg.Hook.Sysfs = &a
	}
	if *cfg.Hook.HooksPath == "" {
		a := os.Getenv(hooks.EnvHooks)
		cfg.Hook.HooksPath = &a
	}
	if *cfg.Hook.HooksPath == "" {
		cfg.Hook.HooksPath = &hooks.DefaultHookPath
	}

	if *cfg.Pod.Cri.Endpoint == "" {
		a := os.Getenv(pods.EnvCri)
		cfg.Pod.Cri.Endpoint = &a
	}
	if *cfg.Pod.ForceProc == false {
		if _, ok := os.LookupEnv(pods.EnvForceProcfs); ok {
			a := true
			cfg.Pod.ForceProc = &a
		}
	}
	if *cfg.Pod.Cri.PodName == "" {
		a := os.Getenv(pods.EnvPodName)
		cfg.Pod.Cri.PodName = &a
	}

	cfg.Logger.Name = appName
	if *cfg.Logger.JaegerEndpoint == "" {
		a := os.Getenv(logger.EnvLoggerJaegerEndpoint)
		cfg.Logger.JaegerEndpoint = &a
	}

	if len(runPathsArg) == 0 {
		runPathsArg.Set(os.Getenv(pods.EnvRunPaths))
	}
	cfg.Pod.Cri.RunPaths = runPathsArg

	return &cfg, flApiAddr
}

func main() {
	var (
		t   *trace.Tracer
		err error
	)

	flag.Usage = usage

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, addr := getConfig()
	if t, err = trace.NewTracer(ctx, cfg); err != nil {
		log.Fatal("Failed to create new tracer: ", err)
		return
	}
	defer t.Destroy()

	router := api.NewRouter(t)

	log.Printf("Listening for incoming API requests at %s", *addr)

	if err = router.Run(*addr); err != nil {
		log.Fatal("Failed to run the server:", err)
	}
}
