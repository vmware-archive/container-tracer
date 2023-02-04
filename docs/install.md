# Installation

Container-tracer supports three types of installation, depending on the specific use case.

## Installation via Docker
There is a Makefile target for building container-tracer docker images, just run `make docker`.
There is no need to run `make` before that, it compiles everything needed. Two ready to use Docker
images are produced, bundled with all dependencies:  
`vmware-labs/container-tracer/tracer-node`  
`vmware-labs/container-tracer/tracer-svc`  
There are different make targets for each of them, so they can be built independently:  
`make docker_tracer` builds `vmware-labs/container-tracer/tracer-node`  
`make docker_service` builds `vmware-labs/container-tracer/tracer-svc`  
The `vmware-labs/container-tracer/tracer-node` image can be used for a local container-tracer setup,
if it is not part of a cluster:
- Run a privileged container using `vmware-labs/container-tracer/tracer-node` image.  
- If everything is ok, the container port `:8080` is exposed and can be used to interact with the tracer,
  using the [REST API](container-tracer-api.md).

## Installation on Kubernetes
Kubernetes is the primary target of container-tracer. There are several installation variants, supported
by Kustomize, but before that there are three important steps:  
- Set the docker repository for the images. As container-tracer is still in its early development
  stage, the images are not yet optimized. That's why no default docker repository is configured:  
    -  Point the `DOCKER_REPO` variable in the top `Makefile` to your docker repository:  
       `DOCKER_REPO=my.docker.repo/`  
    - Prefix the both images in `install/base/kustomization.yaml` with your docker repository:  
       `newName: my.docker.repo/vmware-labs/container-tracer/tracer-svc`  
       `newName: my.docker.repo/vmware-labs/container-tracer/tracer-node`  
- Build the docker images `make docker`  
- Push the images into the repository:  
  `docker push my.docker.repo/vmware-labs/container-tracer/tracer-svc:latest`  
  `docker push my.docker.repo/vmware-labs/container-tracer/tracer-node:latest`  


When the images are in your docker repository, use one of the following installation options:
- Base - serves the container-tracer API over HTTP:  
Run `kubectl apply -k install/base/`.  
- TLS with self-signed certificate - serves the container-tracer API over HTTPS, using a self-signed
certificate. The certificate is managed by the [cert-manager](https://cert-manager.io/docs/installation/)
Kubernetes component, therefore it has to be installed on the cluster.  
Run `kubectl apply -k install/tls/self-signed/`.  
- TLS with custom certificate - serves the container-tracer API over HTTPS, using external certificate.
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
- Build `cmd/tracer-node/tracer-node` and copy the `tracer-node` binary to the desired installation
  location.  
- Copy the `trace-hooks` directory to the desired installation location.
- Install [trace-cruncher](https://github.com/vmware/trace-cruncher) and all its dependencies.  
- Run `tracer-node` with root privileges. It needs the `trace-hooks` directory and searches for it
  in the current directory. You can specify its location using the `TRACER_HOOKS` environment variable
  or the `--trace-hooks` argument:  
  `tracer-node --trace-hooks <path to the trace-hooks directory>`  
- If everything is ok, it will print the REST API endpoint and the available APIs. By default, it listens
  to port `:8080`.  
- That's it. Use the [REST API](container-tracer-api.md) to interact with the tracer. It should
  auto-discover *almost* all containers running on the local system and should be able to run a trace
  session on each of them.
