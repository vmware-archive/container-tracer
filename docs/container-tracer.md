# container-tracer overview
![container-tracer](container-tracer.png)
## Kubernetes
In a Kubernetes environment, container-tracer runs two types of pods:  
- `tracer-node`, one instance on each worker node in the cluster. It is responsible for
   the discovery of containers, running on this node and manage the tracing sessions to
   these containers. This pod runs in a privileged mode and has access to all pods running
   on the node. For a detailed description, see [tracer-node](tracer-node.md).  
- `tracer-svc` - one instance in the entire cluster. It exposes the container-tracer REST API and
   is responsible for broadcasting all API requests to all `trace-node` pods, receiving and
   aggregating the responses. For a detailed description, see [tracer-svc](tracer-svc.md).  
All `tracer-node` pods are bound to a headless ClusterIP service, used for communication with
the `tracer-svc` pod. The `tracer-svc` pod uses a NodePort service to expose the REST API. Look
at [REST API](container-tracer-api.md) for description of the API. Depending on the configuration, each
`trace-node` pod sends the collected traces using Open Telemetry framework to an external database.
Look at [traces flow](container-tracer-flow.md) for description on how the trace collection flow works.
## Standalone
In a standalone or docker installation, only the `tracer-node` is used. It has the same functionality
as in the Kubernetes environment, but with those two differences:  
- Information from the `/proc` filesystem is used to discover the containers running on the system.  
- The REST API is accessible directly, without the `tracer-svc` proxy.