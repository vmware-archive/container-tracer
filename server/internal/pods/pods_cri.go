// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Discover containers running on the local node, using CRI interface.
 */
package pods

import (
	"encoding/json"
	"fmt"
	"time"

	criapi "k8s.io/cri-api/pkg/apis"
	pbuf "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
	ktype "k8s.io/kubernetes/pkg/kubelet/types"
)

var (
	socPrefix       = "unix://"
	defaultRunPaths = []string{
		"/run",
		"/var/run",
	}
	knownCriEndpoints = [...]string{
		"k3s/containerd/containerd.sock",
		"containerd/containerd.sock",
		"cri-dockerd.sock",
		"dockershim.sock",
		"crio/crio.sock",
	}
)

type podCri struct {
	api  criapi.RuntimeService
	podb map[string]*pod
}

type podCriInfo struct {
	Pid int `json:"Pid"`
}

func (p *podCri) criConnect(endpoint *string, runPaths []string) error {
	timeout := 100 * time.Millisecond

	if endpoint != nil && *endpoint != "" {
		if svc, err := remote.NewRemoteRuntimeService(*endpoint, timeout); err != nil {
			return err
		} else {
			p.api = svc
		}
		return nil
	}

	paths := defaultRunPaths
	if len(runPaths) > 0 {
		paths = runPaths
	}

	for _, pt := range paths {
		for _, ep := range knownCriEndpoints {
			sockUrl := socPrefix + pt + "/" + ep
			svc, err := remote.NewRemoteRuntimeService(sockUrl, timeout)
			if err == nil {
				p.api = svc
				print("\nUsing CRI for pods discovery at ", sockUrl, "\n")
				return nil
			}
		}
	}

	return fmt.Errorf("Cannot connect to CRI endpoint")
}

func getCriDiscover(endpoint *string, runPaths []string) (podsDiscover, error) {

	ctr := podCri{
		podb: make(map[string]*pod),
	}

	if err := ctr.criConnect(endpoint, runPaths); err != nil {
		return nil, err
	}

	return ctr, nil
}

func (p *podCri) getPodInfo(cinfo *pbuf.Container, pname *string) error {

	if _, ok := p.podb[*pname]; !ok {
		p.podb[*pname] = &pod{
			Containers: make(map[string]*Container),
		}
	}
	if _, ok := p.podb[*pname].Containers[cinfo.Metadata.Name]; !ok {
		p.podb[*pname].Containers[cinfo.Metadata.Name] = &Container{
			Id:  &cinfo.Metadata.Name,
			Pod: pname,
		}
	}
	cr := p.podb[*pname].Containers[cinfo.Metadata.Name]

	if s, e := p.api.ContainerStatus(cinfo.Id, true); e == nil {
		i := s.GetInfo()
		if v, ok := i["info"]; ok {
			info := podCriInfo{}
			if err := json.Unmarshal([]byte(v), &info); err == nil {
				cr.Tasks = append(cr.Tasks, info.Pid)
			}
		}
	} else {
		return e
	}

	return nil
}

func (p podCri) podScan() (*map[string]*pod, error) {
	// Filter only the running containers
	f := &pbuf.ContainerFilter{
		State: &pbuf.ContainerStateValue{
			State: pbuf.ContainerState_CONTAINER_RUNNING,
		},
	}
	// Get list of all running containers
	r, err := p.api.ListContainers(f)
	if err != nil {
		return nil, err
	}
	// Reset the pods databse
	p.podb = make(map[string]*pod)
	for _, cr := range r {
		if podName, ok := cr.Labels[ktype.KubernetesPodNameLabel]; ok {
			p.getPodInfo(cr, &podName)
		}
	}

	return &p.podb, nil
}
