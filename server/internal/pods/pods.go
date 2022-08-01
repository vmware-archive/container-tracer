// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package condb

import (
	"fmt"
)

type containersDiscover interface {
	contScan() (*map[string]*container, error)
}

type container struct {
	Pids []int `json:"Tasks"`
}

type ContainersDb struct {
	discover   containersDiscover
	containers *map[string]*container
	node       string
}

func getContDiscover(criPath *string, forceProcfs *bool) (containersDiscover, error) {
	var d containersDiscover
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

func NewContainerDb(criPath *string, forceProcfs *bool) (*ContainersDb, error) {

	if d, err := getContDiscover(criPath, forceProcfs); err == nil {
		return &ContainersDb{
			discover: d,
		}, nil
	} else {
		return nil, err
	}
}

func (c *ContainersDb) Scan() error {
	if cdb, err := c.discover.contScan(); err == nil {
		c.containers = cdb
	} else {
		return err
	}

	return nil
}

func (c *ContainersDb) Count() int {
	if c.containers == nil {
		return 0
	}
	return len(*c.containers)
}

func (c *ContainersDb) Print() {
	if c.containers == nil {
		return
	}
	for k, v := range *c.containers {
		fmt.Println(k, *v)
	}
}

func (c *ContainersDb) Get() *map[string]*container {
	return c.containers
}
