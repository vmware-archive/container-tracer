// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Internal in-memory database with all pods and containers running on the local node.
 */
package pods

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	parentPidStr      = "PPid:"
	procfsPathDefault = "/proc"
	EnvForceProcfs    = "TRACER_FORCE_PROCFS"
)

type podsDiscover interface {
	podScan() (*map[string]*pod, error)
}

type PodConfig struct {
	Cri       CriConfig
	ForceProc *bool /* Force using procfs for containers discovery, even if CRI is available. */
}

type Container struct {
	Id, Pod *string
	Parent  []int
	Tasks   []int `json:"Tasks"`
}

type pod struct {
	Containers map[string]*Container
}

type PodDb struct {
	ctx        context.Context
	discover   podsDiscover
	procfsPath *string
	pods       *map[string]*pod
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

func getPodDiscover(ctx context.Context, cfg *PodConfig, procfsPath *string) (podsDiscover, error) {
	var d podsDiscover
	var err error

	if cfg.ForceProc != nil && *cfg.ForceProc {
		return getProcDiscover(ctx, procfsPath)
	}

	if d, err = getCriDiscover(ctx, &cfg.Cri); err == nil {
		return d, err
	} else if d, err = getProcDiscover(ctx, procfsPath); err == nil {
		return d, err
	}

	return nil, err
}

func NewPodDb(ctx context.Context, cfg *PodConfig, procfsPath *string) (*PodDb, error) {

	ppath := procfsPath
	if ppath == nil || *ppath == "" {
		ppath = &procfsPathDefault
	}

	if d, err := getPodDiscover(ctx, cfg, ppath); err == nil {
		db := &PodDb{
			discover:   d,
			ctx:        ctx,
			procfsPath: ppath,
		}
		db.Scan()
		return db, nil
	} else {
		return nil, err
	}
}

func (p *PodDb) getParent(pid int) (int, error) {
	if file, err := os.Open(fmt.Sprintf("%s/%d/status", *p.procfsPath, pid)); err == nil {
		defer file.Close()
		scan := bufio.NewScanner(file)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			words := strings.Fields(scan.Text())
			if len(words) < 2 || words[0] != parentPidStr {
				continue
			}
			if i, err := strconv.Atoi(words[1]); err == nil {
				return i, nil
			} else {
				return 0, err
			}
		}
	} else {
		return 0, err
	}

	return 0, fmt.Errorf("Failed to get the parent")
}

func checkArrayContains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func (p *PodDb) scanParents() {
	for _, pd := range *p.pods {
		for _, cn := range pd.Containers {
			for _, t := range cn.Tasks {
				if ppid, err := p.getParent(t); err == nil {
					if checkArrayContains(cn.Tasks, ppid) {
						continue
					}
					if checkArrayContains(cn.Parent, ppid) {
						continue
					}
					cn.Parent = append(cn.Parent, ppid)
				} else {
					continue
				}
			}
		}
	}
}

func (p *PodDb) Scan() error {
	if cdb, err := p.discover.podScan(); err == nil {
		p.pods = cdb
		p.scanParents()
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
