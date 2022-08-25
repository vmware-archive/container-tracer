// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Internal in-memory database with all pods and containers running on the local node.
 */
package pods

import (
	"fmt"
	"path/filepath"
	"strings"
)

var (
	procfsDefault = "/proc"
)

type podsDiscover interface {
	podScan() (*map[string]*pod, error)
}

type Container struct {
	Id, Pod *string
	Tasks   []int `json:"Tasks"`
}

type pod struct {
	Containers map[string]*Container
}

type PodDb struct {
	discover   podsDiscover
	procfsPath string
	pods       *map[string]*pod
	node       string
}

func hasWildcard(pattern *string) bool {
	if strings.Contains(*pattern, "*") {
		return true
	}
	if strings.Contains(*pattern, "?") {
		return true
	}
	return false
}

func matchName(pattern, name *string) bool {
	if !hasWildcard(pattern) {
		if pattern == name {
			return true
		}
		return false
	}
	if *pattern == "*" {
		return true
	}

	m, _ := filepath.Match(*pattern, *name)
	return m
}

func getContainersFromPod(p *pod, containerName *string) []*Container {

	res := []*Container{}

	if !hasWildcard(containerName) {
		if c, ok := p.Containers[*containerName]; ok {
			res = append(res, c)
		}
		return res
	}

	for cn, c := range p.Containers {
		if !matchName(containerName, &cn) {
			continue
		}
		res = append(res, c)
	}

	return res
}

func (p *PodDb) GetContainers(podName, containerName *string) []*Container {

	res := []*Container{}

	if !hasWildcard(podName) {
		if pd, ok := (*p.pods)[*podName]; ok {
			return getContainersFromPod(pd, containerName)
		}
		return res
	}

	for pn, pd := range *p.pods {
		if !matchName(podName, &pn) {
			continue
		}
		r := getContainersFromPod(pd, containerName)
		res = append(res, r...)
	}

	return res
}

func getPodDiscover(criPath *string, procfsPath *string, forceProcfs *bool) (podsDiscover, error) {
	var d podsDiscover
	var err error

	if forceProcfs != nil && *forceProcfs {
		return getProcDiscover(procfsPath)
	}

	if d, err = getCriDiscover(criPath); err == nil {
		return d, err
	} else if d, err = getProcDiscover(procfsPath); err == nil {
		return d, err
	}

	return nil, err
}

func NewPodDb(criPath *string, forceProcfs *bool) (*PodDb, error) {

	if d, err := getPodDiscover(criPath, &procfsDefault, forceProcfs); err == nil {
		db := &PodDb{
			discover:   d,
			procfsPath: procfsDefault,
		}
		db.Scan()
		return db, nil
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
