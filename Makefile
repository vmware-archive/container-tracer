# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

.PHONY: build clean docker
GO = CGO_ENABLED=0 GOOS=linux go

BIN_TRACER=tracer-node
BIN_SVC=tracer-svc
TRACE_CRUNCER_URL=https://github.com/vmware/trace-cruncher
TRACE_CRUNCER_VER=tracecruncher-v0.4.0

all: build

tracer:
	cd cmd/$(BIN_TRACER) && $(GO) build -o $(BIN_TRACER) .
service:
	cd cmd/$(BIN_SVC) && $(GO) build -o $(BIN_SVC) .

build: tracer service

GIT_SHA=$(shell git rev-parse HEAD)
DOCKER_REPO=
DOCKER_IMAGE=$(DOCKER_REPO)vmware-labs/container-tracer

docker_tracer:
	docker build \
		-f cmd/$(BIN_TRACER)/Dockerfile \
		--build-arg TRACE_CRUNCER_URL=${TRACE_CRUNCER_URL} \
		--build-arg TRACE_CRUNCER_VER=${TRACE_CRUNCER_VER} \
		--label "git_sha=$(GIT_SHA)" \
		-t $(DOCKER_IMAGE)/$(BIN_TRACER):$(GIT_SHA) \
		-t $(DOCKER_IMAGE)/$(BIN_TRACER):latest \
		.

docker_service:
	docker build \
		-f cmd/$(BIN_SVC)/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t $(DOCKER_IMAGE)/$(BIN_SVC):$(GIT_SHA) \
		-t $(DOCKER_IMAGE)/$(BIN_SVC):latest \
		.

docker: docker_tracer docker_service

clean:
	rm -f cmd/$(BIN_TRACER)/$(BIN_TRACER) \
	rm -f cmd/$(BIN_SVC)/$(BIN_SVC)
