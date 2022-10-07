// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Implementation of the common tracer context, used to tie together all su
 */
package tracerctx

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/pods"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracehook"
)

type Tracer struct {
	pods     *pods.PodDb
	hooks    *tracehook.TraceHooks
	sessions *sessionDb
}

type TracerConfig struct {
	Cri       *string  /* CRI  endpoint. */
	Procfs    *string  /* /proc fs  mountpoint. */
	Sysfs     *string  /* /sys fs  mountpoint. */
	Hooks     *string  /* Path to directory with trace hooks. */
	RunPaths  []string /* Paths to run directories. */
	ForceProc *bool    /* Force using procfs for containers discovery, even if CRI is available. */
	Verbose   *bool    /* Print informational logs on the standard output. */
}

func NewTracer(cfg *TracerConfig) (*Tracer, error) {
	var err error
	tr := Tracer{}

	rand.Seed(time.Now().Unix())

	if tr.pods, err = pods.NewPodDb(cfg.Cri, cfg.RunPaths, cfg.Procfs, cfg.ForceProc); err != nil {
		return nil, err
	}
	if tr.hooks, err = tracehook.NewTraceHooksDb(cfg.Hooks, cfg.Procfs, cfg.Sysfs); err != nil {
		return nil, err
	}

	if tr.sessions = newSessionDb(); tr.sessions == nil {
		return nil, fmt.Errorf("Failed to create new session database")
	}

	return &tr, nil
}
