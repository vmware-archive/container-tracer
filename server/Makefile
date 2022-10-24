# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

.PHONY: build clean docker
GO = CGO_ENABLED=0 GO111MODULE=on GOOS=linux go

NODE_BINARY=tracer-node

all: build

build:
	cd cmd/$(NODE_BINARY)/ && $(GO) build -o $(NODE_BINARY) .

GIT_SHA=$(shell git rev-parse HEAD)
DOCKER_NODE_IMAGE=vmware/opensource/$(NODE_BINARY)

docker: build docker_build

docker_build:
	docker build \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t $(DOCKER_NODE_IMAGE):$(GIT_SHA) \
		.

clean:
	rm -f cmd/$(NODE_BINARY)/$(NODE_BINARY)
