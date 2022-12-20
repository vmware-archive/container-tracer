# trace-kube overview
![trace-kube](trace-kube.png)
## Kubernetes
In a Kubernetes environment, trace-kube runs two kinds of pods:  
- `tracer-node`, one instance on each worker node of the cluster. It is responsible for
   discovering the containers, running on that node and manage tracing sessions to these containers.
   This pod runs in a privileged mode and has access to all pods running on the node.
   For a detailed description, see [tracer-node](tracer-node.md).  
- `tracer-svc` - one instance in the entire cluster. It exposes the trace-kube REST API and
   is responsible for broadcasting all API requests to all `trace-node` pods, receiving and
   aggregating the replies. For a detailed description, see [tracer-svc](tracer-svc.md).  
All `tracer-node` pods are bound to a headless ClusterIP service, used for communication with
the `tracer-svc` pod. The `tracer-svc` pod uses a NodePort service to expose the REST API. Look
at [REST API](trace-kube-api.md) for description of the API. Depending on the configuration, each
`trace-node` pod sends the collected traces using Open Telemetry framework to an external database.
Look at [traces flow](trace-kube-flow.md) for description on how the trace collection flow works.
## Standalone
In a standalone or docker installation, only the `tracer-node` is used. It has the same functionality
as in the Kubernetes environment, but with these two differences:  
- Information from the `/proc` filesystem is used to discover the containers running on the system.  
- The REST API is accessible directly, without the `tracer-svc` proxy.