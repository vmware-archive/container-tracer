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
	NodeName *string              /* Name of the cluster node */
	Verbose  *bool                /* Print informational logs on the standard output. */
	Hook     tracehook.HookConfig /* User configuration, specific to trace-hooks database */
	Pod      pods.PodConfig       /* User configuration, specific to pods database */
}

func NewTracer(cfg *TracerConfig) (*Tracer, error) {
	var err error
	tr := Tracer{}

	rand.Seed(time.Now().Unix())

	if tr.pods, err = pods.NewPodDb(&cfg.Pod, cfg.Hook.Procfs); err != nil {
		return nil, err
	}
	if tr.hooks, err = tracehook.NewTraceHooksDb(&cfg.Hook); err != nil {
		return nil, err
	}

	if tr.sessions = newSessionDb(); tr.sessions == nil {
		return nil, fmt.Errorf("Failed to create new session database")
	}

	return &tr, nil
}
