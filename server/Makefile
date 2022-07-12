.PHONY: build clean docker
GO = CGO_ENABLED=0 GO111MODULE=on go

BINARY=trace-kube

all: build

build:
	cd cmd/tracer/ && $(GO) build -o $(BINARY) .

GIT_SHA=$(shell git rev-parse HEAD)

docker: build docker_build

docker_build:
	docker build \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t gitlab.eng.vmware.com/opensource/tracecruncher-api:$(GIT_SHA) \
		.

clean:
	rm -f ./$(BINARY)
