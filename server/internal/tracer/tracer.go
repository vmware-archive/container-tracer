// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package tracer

import (
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/condb"
)

type Tracer struct {
	containers *condb.ContainersDb
}

func NewTracer() (*Tracer, error) {
	var (
		err error
		c   *condb.ContainersDb
	)
	if c, err = condb.NewContainerDb(); err != nil {
		return nil, err
	}

	return &Tracer{
		containers: c,
	}, nil
}
