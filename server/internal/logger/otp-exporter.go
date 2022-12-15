// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Implementation of Open Telemetry trace exporter
 */
package logger

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/exporters/jaeger"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	jaegerDefaultService = "jaeger-collector"
	jaegerDefaultPort    = int32(14268)
	jaegerDefaultRoute   = "api/traces"
)

func jaegerAutoDetect(ctx context.Context) (*string, error) {
	var kclient *k8s.Clientset
	svcFilter := meta.ListOptions{
		FieldSelector: "metadata.name=" + jaegerDefaultService,
	}

	if kconfig, err := rest.InClusterConfig(); err != nil {
		return nil, err
	} else if kclient, err = k8s.NewForConfig(kconfig); err != nil {
		return nil, err
	}

	if ep, err := kclient.CoreV1().Endpoints("").List(ctx, svcFilter); err == nil {
		for _, e := range ep.Items {
			for _, sub := range e.Subsets {
				for _, p := range sub.Ports {
					if p.Port != jaegerDefaultPort {
						continue
					}
					url := fmt.Sprintf("http://%s:%d/%s", jaegerDefaultService, p.Port, jaegerDefaultRoute)
					return &url, nil
				}
			}
		}
	} else {
		return nil, err
	}

	return nil, fmt.Errorf("Cannot find default ", jaegerDefaultService, " on port ", jaegerDefaultPort)
}

func jaegerExporter(ctx context.Context, endpoint *string) (*sdk.SpanExporter, error) {
	var res sdk.SpanExporter
	var err error
	url := endpoint

	if url == nil || *url == "" || *url == "auto" {
		if url, err = jaegerAutoDetect(ctx); err != nil {
			return nil, err
		}
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(*url)))
	if err == nil {
		log.Printf("Connected to jaeger at %s", *url)
		res = sdk.SpanExporter(exp)
		return &res, err
	}

	return nil, err
}
