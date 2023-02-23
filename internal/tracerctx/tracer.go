// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Implementation of the common tracer context, used to tie together all su
 */
package tracerctx

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"

	"github.com/vmware-labs/container-tracer/internal/logger"
	"github.com/vmware-labs/container-tracer/internal/pods"
	"github.com/vmware-labs/container-tracer/internal/tracehook"
)

type Tracer struct {
	pods     *pods.PodDb
	hooks    *tracehook.TraceHooks
	sessions *sessionDb
	logger   *logger.Logger
	node     *string
}

type TracerConfig struct {
	NodeName *string              /* Name of the cluster node */
	Verbose  *bool                /* Print informational logs on the standard output. */
	Hook     tracehook.HookConfig /* User configuration, specific to trace-hooks database */
	Pod      pods.PodConfig       /* User configuration, specific to pods database */
	Logger   logger.LoggerConfig  /* User configuration, specific to trace logger */
}

func setRandomSeed(nodeName *string) {
	h := fnv.New64a()
	h.Write([]byte(*nodeName))
	h.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	rand.Seed(int64(h.Sum64()))
}

func NewTracer(ctx context.Context, cfg *TracerConfig) (*Tracer, error) {
	var err error
	tr := Tracer{
		node: cfg.NodeName,
	}

	setRandomSeed(cfg.NodeName)

	if tr.pods, err = pods.NewPodDb(ctx, &cfg.Pod, cfg.Hook.Procfs); err != nil {
		return nil, err
	}
	if tr.hooks, err = tracehook.NewTraceHooksDb(&cfg.Hook); err != nil {
		return nil, err
	}

	if tr.sessions = newSessionDb(); tr.sessions == nil {
		return nil, fmt.Errorf("Failed to create new session database")
	}

	if tr.logger, err = logger.NewLogger(ctx, &cfg.Logger); err != nil {
		return nil, err
	}

	return &tr, nil
}

func (t *Tracer) Destroy() {
	t.logger.Destroy()
}
