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
		"containerd/containerd.sock",
		"cri-dockerd.sock",
		"dockershim.sock",
		"crio/crio.sock",
		"k3s/containerd/containerd.sock",
	}

	EnvCri      = "TRACER_CRI_ENDPOINT"
	EnvRunPaths = "TRACER_RUN_PATHS"
	EnvPodName  = "TRACER_POD_NAME"
)

type CriConfig struct {
	Endpoint *string  /* CRI endpoint. */
	RunPaths []string /* Paths to run directories. */
	PodName  *string  /* Name of the tracer pod */
}

type podCri struct {
	api  criapi.RuntimeService
	podb map[string]*pod
}

type podCriInfo struct {
	Pid int `json:"Pid"`
}

/* Verify if tracer pod is part of this CRI database */
func criVerify(podName *string, api *criapi.RuntimeService) bool {
	if podName == nil || *podName == "" || api == nil {
		return true
	}

	// Filter only the running containers
	f := &pbuf.ContainerFilter{
		State: &pbuf.ContainerStateValue{
			State: pbuf.ContainerState_CONTAINER_RUNNING,
		},
	}

	// Get list of all running containers
	r, err := (*api).ListContainers(f)
	if err != nil {
		return false
	}

	// Look for tracer pod into the database
	for _, cr := range r {
		if pName, ok := cr.Labels[ktype.KubernetesPodNameLabel]; ok {
			if pName == *podName {
				return true
			}
		}
	}

	return false
}

func (p *podCri) criConnect(cfg *CriConfig) error {
	timeout := 100 * time.Millisecond

	if cfg.Endpoint != nil && *cfg.Endpoint != "" {
		if svc, err := remote.NewRemoteRuntimeService(*cfg.Endpoint, timeout); err != nil {
			return err
		} else {
			p.api = svc
		}
		return nil
	}

	paths := defaultRunPaths
	if len(cfg.RunPaths) > 0 {
		paths = cfg.RunPaths
	}

	for _, pt := range paths {
		for _, ep := range knownCriEndpoints {
			sockUrl := socPrefix + pt + "/" + ep
			svc, err := remote.NewRemoteRuntimeService(sockUrl, timeout)
			if err == nil && criVerify(cfg.PodName, &svc) {
				p.api = svc
				print("\nUsing CRI for pods discovery at ", sockUrl, "\n")
				return nil
			}
		}
	}

	return fmt.Errorf("Cannot connect to CRI endpoint")
}

func getCriDiscover(cfg *CriConfig) (podsDiscover, error) {

	ctr := podCri{
		podb: make(map[string]*pod),
	}

	if err := ctr.criConnect(cfg); err != nil {
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
