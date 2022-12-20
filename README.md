# trace-kube
## Overview
The trace-kube project brings the power of the Linux kernel tracing to Kubernetes. It leverages
existing kernel tracing frameworks such as ftrace, perf, ebpf to trace workloads running on
a Kubernetes cluster. Designed as a native Kubernetes application, its main goal is to be simple
and efficient in doing one thing - collecting low level system traces per container.

## Try it out
### Prerequisites
- Linux kernel with enabled ftrace, perf and ebpf, or at least one of them.
  Almost all of the kernels, shipped with major Linux distributions meet that requirement.  
- Open Telemetry and Jaeger installed on the system / cluster. Although this is not a mandatory
  requirement, it is a good to have. Trace-kube does not store collected traces. All it can do is to
  dump them on the console, or send them to an external database using Open Telemetry.  
- Root privileges on the system / cluster.

### Build
Trace-kube uses Makefile for building, so just type `make` in the top directory of the project.
By default, it builds two applications:  
`cmd/tracer-node/tracer-node`  
`cmd/tracer-svc/tracer-svc`  
There are different make targets for each of them, so they can be compiled independently:  
`make tracer` compiles `cmd/tracer-node/tracer-node`  
`make service` compiles `cmd/tracer-svc/tracer-svc`  

### Standalone installation
There is no Makefile target for a standalone installation, please use **Installation via Docker**
or **Installation on Kubernetes**. However, you can install it by hand:
- Build `cmd/tracer-node/tracer-node` and copy `tracer-node` binary to desired installation location.  
- Copy `trace-hooks` directory to desired installation location.
- Install [trace-cruncher](https://github.com/vmware/trace-cruncher) and all its dependencies.  
- Run `tracer-node` with root privileges. It needs `trace-hooks` directory and by default looks for it
  in current directory. You can specify its location using the `TRACER_HOOKS` environment variable or
  `--trace-hooks` argument:  
  `tracer-node --trace-hooks <path to trace-hooks directory>`  
- If everything is ok, it will print the REST API endpoint and available APIs. By default, it listens
  to port `:8080`.  
- That's it. Use the REST API to interact with the tracer. It should auto-discover *almost* all containers
  running on the local system and should be able to run trace session on each of them.

### Installation via Docker
There is a Makefile target to build trace-kube docker images, just run `make docker`. There is no need
to run `make` before that, it compiles everything needed. Two ready for use Docker images are produced,
bundled with all dependencies:  
`vmware/opensource/trace-kube/tracer-node`  
`vmware/opensource/trace-kube/tracer-svc`  
There are different make targets for each of them, so they can be build independently:  
`make docker_tracer` builds `vmware/opensource/trace-kube/tracer-node`  
`make docker_service` builds `vmware/opensource/trace-kube/tracer-svc`  
The `vmware/opensource/trace-kube/tracer-node` can be used for a trace-kube local installation, if
it is not part of a cluster:
- Run a privileged container using `vmware/opensource/trace-kube/tracer-node` image.  
- If everything is ok, the container port `:8080` is exposed and can be used to interact with the tracer,
  using the REST API.

### Installation on Kubernetes
Kubernetes is the primary target for trace-kube, There is a Kubernetes deployment file `trace-kube.yaml`,
but before that there are three important steps:  
- Set the docker repository for the docker images. As trace-kube is still in its early development
  stage, the images are not optimized yet. That's why no default docker repository is configured:  
    -  Point the `DOCKER_REPO` variable in the top `Makefile` to you docker repository:  
       `DOCKER_REPO=my.docker.repo/`  
    - Prefix the both images in `trace-kube.yaml` with your docker repository:  
       `image: my.docker.repo/vmware/opensource/trace-kube/tracer-svc:latest`  
       `image: my.docker.repo/vmware/opensource/trace-kube/tracer-node:latest`  
- Build the docker images `make docker`  
- Push the built images to the repository:  
  `docker push my.docker.repo/vmware/opensource/trace-kube/tracer-svc:latest`  
  `docker push my.docker.repo/vmware/opensource/trace-kube/tracer-node:latest`  
- Deploy `trace-kube` on the Kubernetes cluster:  
  `kubectl apply -f trace-kube.yaml`  
- If everything is ok, there should be `tracer-node` pods running on each Kubernetes node and
  a `tracer-svc` pod, which serves the REST API.

## Documentation
Look at the [trace-kube documentation](docs) for a detailed explanation of the trace-kube architecture
and description of the REST API.  
Index of available documentation:
- [trace-kube overview](docs/trace-kube.md)
- [trace-kube-api](docs/trace-kube-api.md)
- [trace-kube-flow](docs/trace-kube-flow.md)
- [tracer-node](docs/tracer-node.md)
- [tracer-svc](docs/tracer-svc.md)
- [trace-hooks desription](docs/trace-hooks.md)
- [ftrace hooks](trace-hooks/ftrace/README.md)

## Contributing
The trace-kube project team welcomes contributions from the community. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Trace-kube is available under the [GPLv2.0 or later license](LICENSE).