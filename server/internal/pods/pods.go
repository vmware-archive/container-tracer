// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package pods

import (
	"fmt"
)

type podsDiscover interface {
	podScan() (*map[string]*pod, error)
}

type container struct {
	Tasks []int `json:"Tasks"`
}

type pod struct {
	Containers map[string]*container
}

type PodDb struct {
	discover podsDiscover
	pods     *map[string]*pod
	node     string
}

func getPodDiscover(criPath *string, forceProcfs *bool) (podsDiscover, error) {
	var d podsDiscover
	var err error

	if forceProcfs != nil && *forceProcfs {
		return getProcDiscover()
	}

	if d, err = getCriDiscover(criPath); err == nil {
		return d, err
	} else if d, err = getProcDiscover(); err == nil {
		return d, err
	}

	return nil, err
}

func NewPodDb(criPath *string, forceProcfs *bool) (*PodDb, error) {

	if d, err := getPodDiscover(criPath, forceProcfs); err == nil {
		return &PodDb{
			discover: d,
		}, nil
	} else {
		return nil, err
	}
}

func (p *PodDb) Scan() error {
	if cdb, err := p.discover.podScan(); err == nil {
		p.pods = cdb
	} else {
		return err
	}

	return nil
}

func (p *PodDb) Count() int {
	if p.pods == nil {
		return 0
	}
	return len(*p.pods)
}

func (p *PodDb) Print() {
	if p.pods == nil {
		return
	}
	for k, v := range *p.pods {
		fmt.Println(k, *v)
	}
}

func (p *PodDb) Get() *map[string]*pod {
	return p.pods
}
