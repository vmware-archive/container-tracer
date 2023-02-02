# Installation

Container-tracer supports three types of installation, depending on the specific use case.

## Installation via Docker
There is a Makefile target to build container-tracer docker images, just run `make docker`. There is no need
to run `make` before that, it compiles everything needed. Two ready for use Docker images are produced,
bundled with all dependencies:  
`vmware-labs/container-tracer/tracer-node`  
`vmware-labs/container-tracer/tracer-svc`  
There are different make targets for each of them, so they can be build independently:  
`make docker_tracer` builds `vmware-labs/container-tracer/tracer-node`  
`make docker_service` builds `vmware-labs/container-tracer/tracer-svc`  
The `vmware-labs/container-tracer/tracer-node` can be used for a container-tracer local installation, if
it is not part of a cluster:
- Run a privileged container using `vmware-labs/container-tracer/tracer-node` image.  
- If everything is ok, the container port `:8080` is exposed and can be used to interact with the tracer,
  using the [REST API](container-tracer-api.md).

## Installation on Kubernetes
Kubernetes is the primary target for container-tracer. There are several installation variants, supported
by Kustomize, but before that there are three important steps:  
- Set the docker repository for the docker images. As container-tracer is still in its early development
  stage, the images are not optimized yet. That's why no default docker repository is configured:  
    -  Point the `DOCKER_REPO` variable in the top `Makefile` to you docker repository:  
       `DOCKER_REPO=my.docker.repo/`  
    - Prefix the both images in `install/base/kustomization.yaml` with your docker repository:  
       `newName: my.docker.repo/vmware-labs/container-tracer/tracer-svc`  
       `newName: my.docker.repo/vmware-labs/container-tracer/tracer-node`  
- Build the docker images `make docker`  
- Push the built images to the repository:  
  `docker push my.docker.repo/vmware-labs/container-tracer/tracer-svc:latest`  
  `docker push my.docker.repo/vmware-labs/container-tracer/tracer-node:latest`  


When the images are in your docker repository, use one of these installation options:
- Base - serve the container-tracer API on HTTP:  
Run `kubectl apply -k install/base/`.  
- TLS with self-signed certificate - serve the container-tracer API on HTTPS, using self-signed
certificate. The certificate is managed by the [cert-manager](https://cert-manager.io/docs/installation/)
Kubernetes component, so it must be installed on the cluster.  
Run `kubectl apply -k install/tls/self-signed/`.  
- TLS with custom certificate - serve the container-tracer API on HTTPS, using external certificate.
If you have your own trusted certificate, it can be used by container-tracer. Modify the
`install/tls/external/kustomization.yaml` file with locations of your certificate and key files:  
  `- tls.crt=/path/to/your/tls.cert`  
  `- tls.key=/path/to/your/tls.key`  
and run `kubectl apply -k install/tls/external/`.  


If everything is ok, there should be `tracer-node` pods running on each Kubernetes node and
a `tracer-svc` pod, which serves the [REST API](container-tracer-api.md).

## Standalone installation
There is no Makefile target for a standalone installation, please use **Installation via Docker**
or **Installation on Kubernetes**. However, you can install it by hand:
- Build `cmd/tracer-node/tracer-node` and copy `tracer-node` binary to desired installation location.  
- Copy `trace-hooks` directory to desired installation location.
- Install [trace-cruncher](https://github.com/vmware/trace-cruncher) and all its dependencies.  
- Run `tracer-node` with root privileges. It needs `trace-hooks` directory and by default looks for it
  in the current directory. You can specify its location using the `TRACER_HOOKS` environment variable or
  `--trace-hooks` argument:  
  `tracer-node --trace-hooks <path to the trace-hooks directory>`  
- If everything is ok, it will print the REST API endpoints and available APIs. By default, it listens
  to port `:8080`.  
- That's it. Use the [REST API](container-tracer-api.md) to interact with the tracer. It should
  auto-discover *almost* all containers running on the local system and should be able to run trace
  session on each of them.
