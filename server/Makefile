# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

.PHONY: build clean docker
GO = CGO_ENABLED=0 GO111MODULE=on GOOS=linux go

BIN_TRACER=tracer-node

all: build

tracer:
	cd cmd/$(BIN_TRACER) && $(GO) build -o $(BIN_TRACER) .

build: tracer

GIT_SHA=$(shell git rev-parse HEAD)
DOCKER_NODE_IMAGE=vmware/opensource

docker_tracer:
	docker build \
		-f cmd/$(BIN_TRACER)/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t $(DOCKER_NODE_IMAGE)/$(BIN_TRACER):$(GIT_SHA) \
		-t $(DOCKER_NODE_IMAGE)/$(BIN_TRACER):latest \
		.
docker: docker_tracer

clean:
	rm -f cmd/$(BIN_TRACER)/$(BIN_TRACER)

