// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package tracerctx

import (
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/condb"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracehook"
)

type Tracer struct {
	containers *condb.ContainersDb
	hooks      *tracehook.TraceHooks
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
		c   *condb.ContainersDb
		t   *tracehook.TraceHooks
	)

	if c, err = condb.NewContainerDb(cfg.Cri, cfg.ForceProc); err != nil {
		return nil, err
	}
	if t, err = tracehook.NewTraceHooksDb(cfg.Hooks); err != nil {
		return nil, err
	}

	return &Tracer{
		containers: c,
		hooks:      t,
	}, nil
}
