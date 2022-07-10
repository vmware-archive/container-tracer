// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package tracer

type Tracer struct {
}

func NewTracer() (*Tracer, error) {
	return &Tracer{}, nil
}
