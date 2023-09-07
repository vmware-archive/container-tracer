// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Implementation of Open Telemetry trace exporter
 */
package logger

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	logger "go.opentelemetry.io/otel/trace"
)

var (
	loogerCloseTimeout      = time.Second * 5
	EnvLoggerJaegerEndpoint = "TRACER_JEAGER_ENDPOINT"
)

type LogJob struct {
	Name    string
	File    string
	Node    string
	Pod     string
	Job     string
	Session string
}

type LoggerConfig struct {
	JaegerEndpoint *string
	Name           string
}

type logWorker struct {
	log    *LogJob
	span   logger.Span
	ctx    context.Context
	cancel context.CancelFunc
	count  int
}

type Logger struct {
	ctx        context.Context
	provider   *sdk.TracerProvider
	tracer     logger.Tracer
	logWorkers map[string]*logWorker
}

func (l *Logger) Destroy() {
	ctx, cancel := context.WithTimeout(l.ctx, loogerCloseTimeout)
	defer cancel()
	l.provider.Shutdown(ctx)
}

func readLine(f *os.File) (*[]byte, error) {
	r := bufio.NewReader(f)
	if line, p, err := r.ReadLine(); err == nil {
		if p == false {
			return &line, nil
		}
		ln := make([]byte, len(line))
		copy(ln, line)
		for {
			line, p, err := r.ReadLine()
			if err != nil {
				return nil, err
			}
			ln = append(ln, line...)
			if p == false {
				return &ln, nil
			}
		}
	} else {
		return nil, err
	}
}

func (l *Logger) readFile(job *logWorker) error {
	f, err := os.Open(job.log.File)
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		line, err := readLine(f)
		if err != nil {
			return err
		}
		select {
		case <-job.ctx.Done():
			return job.ctx.Err()
		default:
			_, sp := l.tracer.Start(job.ctx, "trace")
			sp.AddEvent(string(*line))
			sp.End()
			job.count++
		}
	}
	return fmt.Errorf("Completed reading file", job.log.File)
}

func (l *Logger) delCompleted() {
	for f, w := range l.logWorkers {
		if w.ctx.Err() != nil {
			w.span.End()
			log.Printf("Completed trace job %s: %d traces collected", w.log.Name, w.count)
			delete(l.logWorkers, f)
		}
	}
}

func (l *Logger) RunLogJob(log *LogJob) error {
	l.delCompleted()
	if _, ok := l.logWorkers[log.File]; ok {
		return nil
	}

	ctx, cancel := context.WithCancel(l.ctx)
	ctxp, span := l.tracer.Start(ctx, log.Name)
	span.SetAttributes(attribute.Key("node").String(log.Node))
	span.SetAttributes(attribute.Key("pod").String(log.Pod))
	span.SetAttributes(attribute.Key("traceJob").String(log.Job))
	span.SetAttributes(attribute.Key("traceSession").String(log.Session))

	l.logWorkers[log.File] = &logWorker{
		log:    log,
		cancel: cancel,
		ctx:    ctxp,
		span:   span,
	}

	span.AddEvent(log.Name)
	go l.readFile(l.logWorkers[log.File])

	return nil
}

func (l *Logger) StopLogJob(log *LogJob) error {

	if w, ok := l.logWorkers[log.File]; ok {
		w.cancel()
		l.delCompleted()
		return nil
	}
	return fmt.Errorf("No log job for ", log.File)
}

func NewLogger(ctx context.Context, cfg *LoggerConfig) (*Logger, error) {
	var exp *sdk.SpanExporter
	var err error

	l := Logger{
		ctx:        ctx,
		logWorkers: make(map[string]*logWorker),
	}

	if cfg.JaegerEndpoint != nil && *cfg.JaegerEndpoint != "" {
		exp, err = jaegerExporter(ctx, cfg.JaegerEndpoint)
	}

	if err != nil {
		return nil, err
	}

	l.provider = sdk.NewTracerProvider(
		sdk.WithBatcher(*exp),
		sdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Name),
		)),
	)
	if l.provider == nil {
		return nil, fmt.Errorf("Failed to init a trace provider")
	}

	otel.SetTracerProvider(logger.TracerProvider(l.provider))
	l.tracer = l.provider.Tracer(cfg.Name)

	return &l, nil
}
