# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

.PHONY: build clean docker
GO = CGO_ENABLED=0 GO111MODULE=on GOOS=linux go

BIN_TRACER=tracer-node
BIN_SVC=tracer-svc

all: build

tracer:
	cd cmd/$(BIN_TRACER) && $(GO) build -o $(BIN_TRACER) .
service:
	cd cmd/$(BIN_SVC) && $(GO) build -o $(BIN_SVC) .

build: tracer service

GIT_SHA=$(shell git rev-parse HEAD)
DOCKER_NODE_IMAGE=vmware/opensource/trace-kube

docker_tracer:
	docker build \
		-f cmd/$(BIN_TRACER)/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t $(DOCKER_NODE_IMAGE)/$(BIN_TRACER):$(GIT_SHA) \
		-t $(DOCKER_NODE_IMAGE)/$(BIN_TRACER):latest \
		.

docker_service:
	docker build \
		-f cmd/$(BIN_SVC)/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t $(DOCKER_NODE_IMAGE)/$(BIN_SVC):$(GIT_SHA) \
		-t $(DOCKER_NODE_IMAGE)/$(BIN_SVC):latest \
		.

docker: docker_tracer docker_service

clean:
	rm -f cmd/$(BIN_TRACER)/$(BIN_TRACER) \
	rm -f cmd/$(BIN_SVC)/$(BIN_SVC)
