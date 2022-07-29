// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package condb

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

type containerCri struct {
	api   criapi.RuntimeService
	pids  []int
	condb map[string]*container
}

type containerCriInfo struct {
	Pid int `json:"Pid"`
}

func (c *containerCri) criConnect(endpoint *string) error {
	timeout := 100 * time.Millisecond

	if endpoint != nil && *endpoint != "" {
		if svc, err := remote.NewRemoteRuntimeService(*endpoint, timeout); err != nil {
			return err
		} else {
			c.api = svc
		}
		return nil
	}

	for _, ep := range knownCriEndpoints {
		svc, err := remote.NewRemoteRuntimeService(ep, timeout)
		if err == nil {
			c.api = svc
			return nil
		}
	}

	return fmt.Errorf("Cannot connect to CRI endpoint")
}

func getCriDiscover(endpoint *string) (containersDiscover, error) {

	ctr := containerCri{
		pids:  make([]int, 0),
		condb: make(map[string]*container),
	}

	if err := ctr.criConnect(endpoint); err != nil {
		return nil, err
	}

	return ctr, nil
}

func (c *containerCri) getContainerInfo(id, pod string) error {

	if _, ok := c.condb[pod]; !ok {
		c.condb[pod] = &container{}
	}

	if s, e := c.api.ContainerStatus(id, true); e == nil {
		i := s.GetInfo()
		if v, ok := i["info"]; ok {
			info := containerCriInfo{}
			if err := json.Unmarshal([]byte(v), &info); err == nil {
				c.condb[pod].Pids = append(c.condb[pod].Pids, info.Pid)
			}
		}
	} else {
		return e
	}

	return nil
}

func (c containerCri) contScan() (*map[string]*container, error) {
	// Filter only the running containers
	f := &pbuf.ContainerFilter{
		State: &pbuf.ContainerStateValue{
			State: pbuf.ContainerState_CONTAINER_RUNNING,
		},
	}
	// Get list of all running containers
	r, err := c.api.ListContainers(f)
	if err != nil {
		return nil, err
	}
	// Reset the containers databse
	c.condb = make(map[string]*container)
	for _, cr := range r {
		if podName, ok := cr.Labels[ktype.KubernetesPodNameLabel]; ok {
			c.getContainerInfo(cr.Id, podName)
		}
	}

	return &c.condb, nil
}
