// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package tracerctx

import (
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/pods"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracehook"
)

type Tracer struct {
	pods  *pods.PodDb
	hooks *tracehook.TraceHooks
}

type TracerConfig struct {
	Cri       *string /* CRI  endpoint. */
	Hooks     *string /* Path to directory with trace hooks. */
	ForceProc *bool   /* Force using procfs for containers discovery, even if CRI is available. */
	Verbose   *bool   /* Print informational logs on the standard output. */
}

func NewTracer(cfg *TracerConfig) (*Tracer, error) {
	var (
		err error
		p   *pods.PodDb
		t   *tracehook.TraceHooks
	)

	if p, err = pods.NewPodDb(cfg.Cri, cfg.ForceProc); err != nil {
		return nil, err
	}
	if t, err = tracehook.NewTraceHooksDb(cfg.Hooks); err != nil {
		return nil, err
	}

	return &Tracer{
		pods:  p,
		hooks: t,
	}, nil
}
