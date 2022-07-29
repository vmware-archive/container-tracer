package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.eng.vmware.com/opensource/tracecruncher-api/api"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracerctx"
)

var (
	envAddress     = "TRACER_API_ADDRESS"
	envCri         = "TRACER_CRI_ENDPOINT"
	envHooks       = "TRACER_HOOKS"
	envForceProcfs = "TRACER_FORCE_PROCFS"
	envVerbose     = "TRACER_VERBOSE"
	defAddress     = ":8080"
	defHooks       = "trace-hooks"
)

func getConfig() (*ctx.TracerConfig, *string) {
	cfg := ctx.TracerConfig{}

	flApiAddr := flag.String("address", "",
		fmt.Sprintf("IP address and port in format IP:port, used for listening for incoming API requests.Can be passed using %s environment variable as well",
			envAddress))
	cfg.Cri = flag.String("cri-endpoint", "",
		fmt.Sprintf("Path to the CRI endpoint. Can be passed using %s environment variable as well.", envCri))
	cfg.Hooks = flag.String("trace-hooks", "",
		fmt.Sprintf("Location of the directory with trace helper applications. Can be passed using %s environment variable as well.", envHooks))
	cfg.ForceProc = flag.Bool("use-procfs", false,
		fmt.Sprintf("Force using procfs for containers discovery, even if CRI is available. Can be passed using %s environment variable as well.", envForceProcfs))
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

	if *cfg.Cri == "" {
		a := os.Getenv(envCri)
		cfg.Cri = &a
	}

	if *cfg.ForceProc == false {
		if _, ok := os.LookupEnv(envForceProcfs); ok {
			a := true
			cfg.ForceProc = &a
		}
	}

	if *cfg.Verbose == false {
		if _, ok := os.LookupEnv(envVerbose); ok {
			a := true
			cfg.Verbose = &a
		}
	}

	if *cfg.Hooks == "" {
		a := os.Getenv(envHooks)
		cfg.Hooks = &a
	}
	if *cfg.Hooks == "" {
		cfg.Hooks = &defHooks
	}

	return &cfg, flApiAddr
}

func main() {
	var (
		t   *ctx.Tracer
		err error
	)

	cfg, addr := getConfig()
	if t, err = ctx.NewTracer(cfg); err != nil {
		return
	}

	router := api.NewRouter(t)

	log.Printf("Listening for incoming API requests at %s", *addr)

	if err = router.Run(*addr); err != nil {
		log.Fatal("Failed to run the server:", err)
	}
}
