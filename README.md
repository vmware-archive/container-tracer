# container-tracer

## Overview
The container-tracer project brings the power of the Linux kernel tracing to Kubernetes. It leverages
existing kernel tracing frameworks such as ftrace, perf, ebpf to trace workloads running on
a Kubernetes cluster. Designed as a native Kubernetes application, its main goal is to be simple
and efficient in doing one thing - collecting low level system traces per container.

## Try it out

### Prerequisites
- Linux kernel with enabled ftrace. Almost all kernels, shipped with major Linux distributions
  meet that requirement.  
- Open Telemetry and Jaeger installed on the system / cluster. Although this is not a mandatory
  requirement, it is a good to have. Container-tracer does not store the collected traces. All it
  can do is to dump them on the console, or send them to an external database using Open Telemetry.  
- Root permissions on the system / cluster.

### Build
Container-tracer uses Makefile to build, so just type `make` in the top directory of the project.
By default, it builds two applications:  
`cmd/tracer-node/tracer-node`  
`cmd/tracer-svc/tracer-svc`  
There are different make targets for each of them, so they can be compiled independently:  
`make tracer` compiles `cmd/tracer-node/tracer-node`  
`make service` compiles `cmd/tracer-svc/tracer-svc`  

### Install

Look at [installation](docs/install.md) for detailed instructions.

### Usage
After installation of the container-tracer, you can interact with it using
a [REST API](docs/container-tracer-api.md).

## Documentation
Look at the [container-tracer documentation](docs) for a detailed explanation of the
container-tracer architecture and a description of the REST API.  
Index of available documentation:
- [installation](docs/install.md)
- [container-tracer overview](docs/container-tracer.md)
- [container-tracer-api](docs/container-tracer-api.md)
- [container-tracer-flow](docs/container-tracer-flow.md)
- [tracer-node](docs/tracer-node.md)
- [tracer-svc](docs/tracer-svc.md)
- [trace-hooks desription](docs/trace-hooks.md)
- [ftrace hooks](trace-hooks/ftrace/README.md)

## Contributing
The container-tracer project team welcomes contributions from the community. For more detailed
information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Container-tracer is available under the [GPLv2.0 or later license](LICENSE).
