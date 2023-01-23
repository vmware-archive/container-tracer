// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Implementation of the container-tracer context, used to tie together all su
 */
package tracekubectx

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	kapi "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	tracerUrlPrefix = "http://"
)

type nodeTracer struct {
	state  kapi.PodPhase
	client *http.Client
	target *url.URL
	ip     string
}

type TraceKube struct {
	config *TraceKubeConfig

	kclient   *k8s.Clientset
	trSync    sync.RWMutex
	tracers   map[string]*nodeTracer
	podFilter meta.ListOptions
	svcFilter meta.ListOptions
}

type TraceKubeConfig struct {
	Verbose     *bool         /* Print informational logs on the standard output. */
	TracersPoll time.Duration /* Polling interval for refreshing the tracers database */
	PodSelector *string       /* Selector for filtering node tracer pods */
	SvcSelector *string       /* Selector for filtering node tracer services */
}

func (t *TraceKube) discoveryTask() {
	tick := time.NewTimer(t.config.TracersPoll)
	<-tick.C
	t.discoverTracers()
}

func (t *TraceKube) updateTracers(ip *string, port int32) {
	for k, n := range t.tracers {
		if n.state == kapi.PodUnknown || n.client != nil {
			continue
		}
		if n.ip == *ip {
			str := tracerUrlPrefix + n.ip + ":" + strconv.FormatInt(int64(port), 10)
			if u, err := url.Parse(str); err == nil {
				fmt.Print("\tAdd node [", k, "] @ ", u, "\n")
				n.client = &http.Client{}
				n.target = u
			}
		}
	}
}

// Auto discover all tracer pods, running on each node from the cluster
func (t *TraceKube) discoverTracers() error {

	t.trSync.Lock()
	defer t.trSync.Unlock()

	// Invalidate all current tracers
	for _, n := range t.tracers {
		n.state = kapi.PodUnknown
	}

	// Get all tracer pods
	if pods, err := t.kclient.CoreV1().Pods("").List(context.TODO(), t.podFilter); err == nil {
		for _, p := range pods.Items {
			if _, ok := t.tracers[p.Name]; !ok {
				t.tracers[p.Name] = &nodeTracer{}
			} else {
				if t.tracers[p.Name].ip != p.Status.PodIP {
					t.tracers[p.Name].client = nil
				}
			}
			t.tracers[p.Name].state = p.Status.Phase
			t.tracers[p.Name].ip = p.Status.PodIP
		}
	} else {
		return err
	}

	newNode := false
	for _, n := range t.tracers {
		if n.state != kapi.PodUnknown && n.client == nil {
			newNode = true
			break
		}
	}

	if newNode == true {
		// Get all tracer endpoints
		if ep, err := t.kclient.CoreV1().Endpoints("").List(context.TODO(), t.svcFilter); err == nil {
			for _, e := range ep.Items {
				for _, sub := range e.Subsets {
					for _, p := range sub.Ports {
						for _, a := range sub.Addresses {
							t.updateTracers(&a.IP, p.Port)
						}
					}
				}
			}
		} else {
			return err
		}
	}

	// Delete invalid tracers
	for n, p := range t.tracers {
		if p.state == kapi.PodUnknown {
			delete(t.tracers, n)
		}
	}

	go t.discoveryTask()
	return nil
}

func NewTraceKube(cfg *TraceKubeConfig) (*TraceKube, error) {
	tk := TraceKube{
		config:  cfg,
		tracers: make(map[string]*nodeTracer),
		podFilter: meta.ListOptions{
			LabelSelector: *cfg.PodSelector,
		},
		svcFilter: meta.ListOptions{
			FieldSelector: *cfg.SvcSelector,
		},
	}

	rand.Seed(time.Now().Unix())

	if kconfig, err := rest.InClusterConfig(); err != nil {
		return nil, err
	} else if tk.kclient, err = k8s.NewForConfig(kconfig); err != nil {
		return nil, err
	}

	if err := tk.discoverTracers(); err != nil {
		return nil, err
	}

	return &tk, nil
}
