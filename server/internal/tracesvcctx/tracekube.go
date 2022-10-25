// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Implementation of the trace-kube context, used to tie together all su
 */
package tracekubectx

import (
	"math/rand"
	"time"
)

type TraceKube struct {
	config *TraceKubeConfig
	//	pods     *pods.PodDb
	//	hooks    *tracehook.TraceHooks
	//	sessions *sessionDb
}

type TraceKubeConfig struct {
	Verbose *bool /* Print informational logs on the standard output. */
}

func NewTraceKube(cfg *TraceKubeConfig) (*TraceKube, error) {
	tk := TraceKube{
		config: cfg,
	}

	rand.Seed(time.Now().Unix())

	return &tk, nil
}
