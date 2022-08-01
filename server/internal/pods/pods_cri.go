// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

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

var knownCriEndpoints = [...]string{
	"unix:///run/containerd/containerd.sock",
	"unix:///var/run/cri-dockerd.sock",
	"unix:///var/run/dockershim.sock",
	"unix:///run/crio/crio.sock",
}

type podCri struct {
	api  criapi.RuntimeService
	pids []int
	podb map[string]*pod
}

type podCriInfo struct {
	Pid int `json:"Pid"`
}

func (p *podCri) criConnect(endpoint *string) error {
	timeout := 100 * time.Millisecond

	if endpoint != nil && *endpoint != "" {
		if svc, err := remote.NewRemoteRuntimeService(*endpoint, timeout); err != nil {
			return err
		} else {
			p.api = svc
		}
		return nil
	}

	for _, ep := range knownCriEndpoints {
		svc, err := remote.NewRemoteRuntimeService(ep, timeout)
		if err == nil {
			p.api = svc
			return nil
		}
	}

	return fmt.Errorf("Cannot connect to CRI endpoint")
}

func getCriDiscover(endpoint *string) (podsDiscover, error) {

	ctr := podCri{
		pids: make([]int, 0),
		podb: make(map[string]*pod),
	}

	if err := ctr.criConnect(endpoint); err != nil {
		return nil, err
	}

	return ctr, nil
}

func (p *podCri) getPodInfo(id, pname string) error {

	if _, ok := p.podb[pname]; !ok {
		p.podb[pname] = &pod{}
	}

	if s, e := p.api.ContainerStatus(id, true); e == nil {
		i := s.GetInfo()
		if v, ok := i["info"]; ok {
			info := podCriInfo{}
			if err := json.Unmarshal([]byte(v), &info); err == nil {
				p.podb[pname].Pids = append(p.podb[pname].Pids, info.Pid)
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
			p.getPodInfo(cr.Id, podName)
		}
	}

	return &p.podb, nil
}
